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
		PrintFieldGraph(function.From[i].Field, "\t")
	}

	for i := 0; i < len(function.To); i++ {
		fmt.Println(function.To[i])
		PrintFieldGraph(function.To[i].Field, "\t")
	}
}

// PrintFieldGraph prints a list of fields with the related fields.
func PrintFieldGraph(field *models.Field, tabs string) {
	fmt.Printf("%v%v\n", tabs, field)

	for i := 0; i < len(field.Fields); i++ {
		if len(field.Fields) != 0 {
			PrintFieldGraph(field.Fields[i], tabs+"\t")
		}
	}
}

// PrintFieldTree prints a tree of fields for a given type to standard output.
func PrintFieldTree(typename string, fields []*models.Field, tabs string) {
	if tabs == "" {
		fmt.Println(tabs + "type " + typename)
	}

	tabs += "\t" // field tab
	for i := 0; i < len(fields); i++ {
		fmt.Println(tabs + fields[i].Name + "\t" + fields[i].Definition)

		if len(fields[i].Fields) != 0 {
			PrintFieldTree(fields[i].Definition, fields[i].Fields, tabs+"\t")
		}
	}
}

// PrintFieldRelation prints the relationship between to and from fields.
func PrintFieldRelation(toFields, fromFields []*models.Field) {
	for i := 0; i < len(toFields); i++ {
		for j := 0; j < len(fromFields); j++ {
			printFieldRelation(toFields[i], fromFields[j])
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
		switch {
		case len(toField.Fields) != 0 && len(fromField.Fields) != 0:
			for i := 0; i < len(toField.Fields); i++ {
				for j := 0; j < len(fromField.Fields); j++ {
					printFieldRelation(toField.Fields[i], fromField.Fields[j])
				}
			}
		case len(toField.Fields) != 0:
			for i := 0; i < len(toField.Fields); i++ {
				printFieldRelation(toField.Fields[i], fromField)
			}
		case len(fromField.Fields) != 0:
			for i := 0; i < len(fromField.Fields); i++ {
				printFieldRelation(toField, fromField.Fields[i])
			}
		default:
			fmt.Printf("%v is not related to %v.\n", toField, fromField)
		}
	}
}

// CountFields returns the number of fields (including subfields) in a field slice.
func CountFields(fields []*models.Field) int {
	if len(fields) == 0 {
		return 0
	}

	for _, field := range fields {
		return 1 + CountFields(field.Fields)
	}

	return 0
}
