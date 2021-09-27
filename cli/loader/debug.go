package loader

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/models"
)

// PrintFieldGraph prints a list of fields with the related fields.
func PrintFieldGraph(fields []*models.Field, tabs string) {
	for _, field := range fields {
		fmt.Println(tabs, field)
		if len(field.Fields) != 0 {
			for _, nestedField := range field.Fields {
				PrintFieldGraph(nestedField.Fields, tabs+"\tNESTED ")
			}
		}
	}
}

// PrintFieldTree prints a tree of fields for a given type to standard output.
func PrintFieldTree(typename string, fields []*models.Field, tabs string) {
	if tabs == "" {
		fmt.Println(tabs + "type " + typename)
	}

	tabs += "\t" // field tab
	for _, field := range fields {
		fmt.Println(tabs + field.Name + "\t" + field.Definition)
		if len(field.Fields) != 0 {
			PrintFieldTree(field.Definition, field.Fields, tabs+"\t")
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

// printFieldRelation prints the relatinship between two fields.
func printFieldRelation(toField *models.Field, fromField *models.Field) {
	if (*toField).From == fromField && (*fromField).To == toField {
		fmt.Printf("To Field %v%v (%v) and From Field %v%v (%v) are related to each other.\n", toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name, fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name)
	} else if (*toField).From == fromField {
		fmt.Printf("To Field %v%v (%v) is related to From Field %v%v (%v).\n", toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name, fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name)
	} else if (*fromField).To == toField {
		fmt.Printf("From Field %v%v (%v) is related to To Field %v%v (%v).\n", fromField.Definition+" ", fromField.Name, fromField.Parent.Package+"."+fromField.Parent.Name, toField.Definition+" ", toField.Name, toField.Parent.Package+"."+toField.Parent.Name)
	} else {
		if len(toField.Fields) != 0 && len(fromField.Fields) != 0 {
			for _, nestedToField := range toField.Fields {
				for _, nestedFromField := range fromField.Fields {
					printFieldRelation(nestedToField, nestedFromField)
				}
			}
		} else if len(toField.Fields) != 0 {
			for _, nestedToField := range toField.Fields {
				printFieldRelation(nestedToField, fromField)
			}
		} else if len(fromField.Fields) != 0 {
			for _, nestedFromField := range fromField.Fields {
				printFieldRelation(toField, nestedFromField)
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
