// DO NOT CHANGE PACKAGE

// Package templates provides a template used by copygen to generate custom code.
package templates

import (
	"github.com/switchupcb/copygen/cli/models"
)

func Generate(gen models.Generator) string {
	content := string(gen.Keep) + "\n"
	for _, function := range gen.Functions {
		content += Function(function) + "\n"
	}
	return content
}

// Function provides generated code for a function.
func Function(function models.Function) string {
	// comment
	fn := generateComment(function) + "\n"

	// signature
	fn += generateSignature(function) + "\n"

	// body
	fn += generateBody(function) + "\n"

	// return
	fn += generateReturn(function) + "\n"

	// end of function
	fn += "}"
	return fn
}

// generateComment generates a function comment.
func generateComment(function models.Function) string {
	var toComment string
	for _, toType := range function.To {
		toComment += toType.Field.Name + ", "
	}
	if len(toComment) != 0 {
		// remove last ", "
		toComment = toComment[:len(toComment)-2]
	}

	var fromComment string
	for _, fromType := range function.From {
		fromComment += fromType.Field.Name + ", "
	}
	if len(fromComment) != 0 {
		// remove last ", "
		fromComment = fromComment[:len(fromComment)-2]
	}

	return "// " + function.Name + " copies a " + fromComment + " to a " + toComment + "."
}

// generateSignature generates a function's signature.
func generateSignature(function models.Function) string {
	sig := "func " + function.Name + "(" + generateParameters(function) + ") {"
	return sig
}

// generateParameters generates the parameters of a function.
func generateParameters(function models.Function) string {
	var parameters string

	// Generate To-Type parameters
	for _, toType := range function.To {
		parameters += toType.Field.VariableName + " "
		parameters += toType.ParameterName() + ", "
	}

	// Generate From-Type parameters
	for _, fromType := range function.From {
		parameters += fromType.Field.VariableName + " "
		parameters += fromType.ParameterName() + ", "
	}

	if len(parameters) == 0 {
		return parameters
	}

	// remove last ", "
	return parameters[:len(parameters)-2]
}

// generateBody generates the body of a function.
func generateBody(function models.Function) string {
	var body string

	// Assign fields to ToType(s)
	for _, toType := range function.To {
		body += "// " + toType.Field.Name + " fields\n"

		for _, toField := range toType.Field.Fields {
			body += toField.FullVariableName("")
			body += " = "
			fromField := toField.From
			if fromField.Options.Convert != "" {
				body += fromField.Options.Convert + "(" + fromField.FullVariableName("") + ")\n"
			} else {
				body += fromField.FullVariableName("") + "\n"
			}
		}
		body += "\n"
	}
	return body
}

// generateReturn generates a return statement for the function.
func generateReturn(function models.Function) string {
	return ""
}
