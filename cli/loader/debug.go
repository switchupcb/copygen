package loader

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/models"
)

// PrintFunctionFields prints all of a functions fields to standard output.
func PrintFunctionFields(function models.Function) {
	for i := 0; i < len(function.To); i++ {
		fmt.Println("TO type " + function.To[i].Package + "." + function.To[i].Name)
		PrintFieldGraph(function.To[i].Fields, "\t")
	}

	for i := 0; i < len(function.From); i++ {
		fmt.Println("FROM type " + function.From[i].Package + "." + function.From[i].Name)
		PrintFieldGraph(function.From[i].Fields, "\t")
	}
}

// PrintFieldGraph prints a list of fields with the related fields.
func PrintFieldGraph(fields []*models.Field, tabs string) {
	for i := 0; i < len(fields); i++ {
		fmt.Printf("%v%v\n", tabs, fields[i])
		if len(fields[i].Fields) != 0 {
			PrintFieldGraph(fields[i].Fields, tabs+"\t")
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
func PrintFieldRelation(toFields []*models.Field, fromFields []*models.Field) {
	for i := 0; i < len(toFields); i++ {
		for j := 0; j < len(fromFields); j++ {
			printFieldRelation(toFields[i], fromFields[j])
		}
	}
}

// printFieldRelation prints the relationship between two fields.
func printFieldRelation(toField *models.Field, fromField *models.Field) {
	if (*toField).From == fromField && (*fromField).To == toField {
		fmt.Printf("To Field %v%v (%v) and From Field %v%v (%v) are related to each other.\n", toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name, fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name)
	} else if (*toField).From == fromField {
		fmt.Printf("To Field %v%v (%v) is related to From Field %v%v (%v).\n", toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name, fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name)
	} else if (*fromField).To == toField {
		fmt.Printf("From Field %v%v (%v) is related to To Field %v%v (%v).\n", fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name, toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name)
	} else {
		if len(toField.Fields) != 0 && len(fromField.Fields) != 0 {
			for i := 0; i < len(toField.Fields); i++ {
				for j := 0; j < len(fromField.Fields); i++ {
					printFieldRelation(toField.Fields[i], fromField.Fields[j])
				}
			}
		} else if len(toField.Fields) != 0 {
			for i := 0; i < len(toField.Fields); i++ {
				printFieldRelation(toField.Fields[i], fromField)
			}
		} else if len(fromField.Fields) != 0 {
			for i := 0; i < len(fromField.Fields); i++ {
				printFieldRelation(toField, fromField.Fields[i])
			}
		} else {
			fmt.Printf("To Field %v%v (%v) is not related to From Field %v%v (%v).\n", toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name, fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name)
		}
	}
}

// CountFields returns the number of fields in a field slice.
func CountFields(fields []*models.Field) int {
	if len(fields) == 0 {
		return 0
	}

	for _, field := range fields {
		return 1 + CountFields(field.Fields)
	}
	return 0
}
