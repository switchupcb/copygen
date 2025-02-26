// DO NOT CHANGE PACKAGE

// Package template provides a template used by copygen to generate custom code.
package template

import (
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// Generate generates code.
// GENERATOR FUNCTION.
// EDITABLE.
// DO NOT REMOVE.
func Generate(gen *models.Generator) (string, error) {
	var content strings.Builder

	content.WriteString(string(gen.Keep) + "\n")
	for i := range gen.Functions {
		content.WriteString(Function(&gen.Functions[i]) + "\n")
	}

	return content.String(), nil
}

// Function provides generated code for a function.
func Function(function *models.Function) string {
	var fn strings.Builder
	fn.WriteString(generateComment(function) + "\n")
	fn.WriteString(generateSignature(function) + "\n")
	fn.WriteString(generateBody(function))
	fn.WriteString(generateReturn(function))
	return fn.String()
}

// generateComment generates a function comment.
func generateComment(function *models.Function) string {
	var toComment strings.Builder
	for i, toType := range function.To {
		if i+1 == len(function.To) {
			toComment.WriteString(toType.Name())
			break
		}

		toComment.WriteString(toType.Name() + ", ")
	}

	var fromComment strings.Builder
	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			fromComment.WriteString(fromType.Name())
			break
		}

		fromComment.WriteString(fromType.Name() + ", ")
	}

	return "// " + function.Name + " copies a " + fromComment.String() + " to a " + toComment.String() + "."
}

// generateSignature generates a function's signature.
func generateSignature(function *models.Function) string {
	return "func " + function.Name + "(" + generateParameters(function) + ") {"
}

// generateParameters generates the parameters of a function.
func generateParameters(function *models.Function) string {
	var parameters strings.Builder
	for _, toType := range function.To {
		parameters.WriteString(toType.Field.VariableName + " " + toType.Name() + ", ")
	}

	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name())
			break
		}

		parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name() + ", ")
	}

	return parameters.String()
}

// generateBody generates the body of a function.
func generateBody(function *models.Function) string {
	var body strings.Builder

	// Assign fields to ToType(s).
	for i, toType := range function.To {
		body.WriteString(generateAssignment(toType))
		if i+1 != len(function.To) {
			body.WriteString("\n")
		}
	}

	return body.String()
}

// generateAssignment generates assignments for a to-type.
func generateAssignment(toType models.Type) string {
	var assign strings.Builder
	assign.WriteString("// " + toType.Name() + " fields\n")

	for _, toField := range toType.Field.AllFields(nil, nil) {
		if toField.From != nil {
			assign.WriteString(toField.FullVariableName("") + " = ")

			fromField := toField.From
			if fromField.Options.Convert != "" {
				assign.WriteString(fromField.Options.Convert + "(" + fromField.FullVariableName("") + ")\n")
			} else if fromField.Options.Cast != "" {
				assign.WriteString(fromField.FullVariableName("") + "." + fromField.Options.Cast + "\n")
			} else {
				switch {
				case toField.FullDefinition() == fromField.FullDefinition():
					assign.WriteString(fromField.FullVariableName("") + "\n")
				case toField.FullDefinition()[1:] == fromField.FullDefinition():
					assign.WriteString("&" + fromField.FullVariableName("") + "\n")
				case toField.FullDefinition() == fromField.FullDefinition()[1:]:
					assign.WriteString("*" + fromField.FullVariableName("") + "\n")
				}
			}
		}
	}

	return assign.String()
}

// generateReturn generates a return statement for the function.
func generateReturn(function *models.Function) string {
	return "}"
}
