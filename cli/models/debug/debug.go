package debug

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/models"
)

// PrintGeneratorFields prints all of a generator's function's fields to standard output.
func PrintGeneratorFields(gen *models.Generator) {
	for i := 0; i < len(gen.Functions); i++ {
		fmt.Println(gen.Functions[i].Name, "{")
		PrintFunctionFields(&gen.Functions[i])
		fmt.Println("}")
		fmt.Println()
	}
}

// PrintFunctionFields prints all of a function's fields to standard output.
func PrintFunctionFields(function *models.Function) {
	for i := 0; i < len(function.From); i++ {
		fmt.Println(function.From[i])
		PrintFieldGraph(function.From[i].Field, "\t", nil)
	}

	for i := 0; i < len(function.To); i++ {
		fmt.Println(function.To[i])
		PrintFieldGraph(function.To[i].Field, "\t", nil)
	}
}

// PrintFieldGraph prints a list of fields with the related fields.
func PrintFieldGraph(field *models.Field, tabs string, cyclic map[*models.Field]bool) {
	if cyclic == nil {
		cyclic = make(map[*models.Field]bool)
	}

	fmt.Printf("%v%v\n", tabs, field)
	cyclic[field] = true
	for _, subfield := range field.Fields {
		if !cyclic[subfield] {
			PrintFieldGraph(subfield, tabs+"\t", cyclic)
		}
	}
}

// PrintFieldTree prints a tree of fields for a given type field to standard output.
func PrintFieldTree(field *models.Field, tabs string, cyclic map[*models.Field]bool) {
	if cyclic == nil {
		cyclic = make(map[*models.Field]bool)
	}

	if tabs == "" {
		fmt.Println(tabs + "type " + field.FullDefinition())
	} else {
		fmt.Println(tabs + field.Name + "\t" + field.FullDefinition())
	}
	cyclic[field] = true

	tabs += "\t" // field tab
	for _, subfield := range field.Fields {
		if !cyclic[subfield] {
			PrintFieldTree(subfield, tabs+"\t", cyclic)
		}
	}
}

// PrintFieldRelation prints the relationship between a list of to and from fields.
func PrintFieldRelation(toFields, fromFields []*models.Field) {
	for _, toField := range toFields {
		for _, fromField := range fromFields {
			printFieldRelation(toField, fromField)
		}
	}
}

// printFieldRelation prints the relationship between two fields.
func printFieldRelation(toField, fromField *models.Field) {
	switch {
	case toField.From == fromField && fromField.To == toField:
		fmt.Printf("%v and %v are related to each other.\n", toField, fromField)
	case toField.From == fromField:
		fmt.Printf("%v is related to %v.\n", toField, fromField)
	case fromField.To == toField:
		fmt.Printf("%v is related to %v.\n", toField, fromField)
	default:
		fmt.Printf("%v is not related to %v.\n", toField, fromField)
	}
}
