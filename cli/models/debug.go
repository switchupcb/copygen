package models

import (
	"fmt"
)

// PrintFunctionFields prints all of a functions fields to standard output.
func PrintFunctionFields(function Function) {
	for i := 0; i < len(function.From); i++ {
		PrintFieldGraph(function.From[i].Field, "\t")
	}
	for i := 0; i < len(function.To); i++ {
		fmt.Println(function.To[i])
		PrintFieldGraph(function.To[i].Field, "\t")
	}
}

// PrintFieldGraph prints a list of fields with the related fields.
func PrintFieldGraph(field *Field, tabs string) {
	fmt.Printf("%v%v\n", tabs, field)
	for i := 0; i < len(field.Fields); i++ {
		if len(field.Fields) != 0 {
			PrintFieldGraph(field.Fields[i], tabs+"\t")
		}
	}
}

// PrintFieldTree prints a tree of fields for a given type to standard output.
func PrintFieldTree(typename string, fields []*Field, tabs string) {
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
func PrintFieldRelation(toFields []*Field, fromFields []*Field) {
	for i := 0; i < len(toFields); i++ {
		for j := 0; j < len(fromFields); j++ {
			printFieldRelation(toFields[i], fromFields[j])
		}
	}
}

// printFieldRelation prints the relationship between two fields.
func printFieldRelation(toField *Field, fromField *Field) {
	if (*toField).From == fromField && (*fromField).To == toField {
		fmt.Printf("To Field %v and From Field %v are related to each other.\n", toField, fromField)
	} else if (*toField).From == fromField {
		fmt.Printf("To Field %v is related to From Field %v.\n", toField, fromField)
	} else if (*fromField).To == toField {
		fmt.Printf("From Field %v is related to To Field %v.\n", toField, fromField)
	} else {
		if len(toField.Fields) != 0 && len(fromField.Fields) != 0 {
			for i := 0; i < len(toField.Fields); i++ {
				for j := 0; j < len(fromField.Fields); i++ { //nolint:staticcheck
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
			fmt.Printf("To Field %v is not related to From Field %v.\n", toField, fromField)
		}
	}
}

// CountFields returns the number of fields in a field slice.
func CountFields(fields []*Field) int {
	if len(fields) == 0 {
		return 0
	}

	for _, field := range fields {
		return 1 + CountFields(field.Fields)
	}
	return 0
}
