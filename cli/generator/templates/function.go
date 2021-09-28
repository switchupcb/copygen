// DO NOT CHANGE PACKAGE
// package templates provides a template used by copygen to generate custom code.
package templates

import (
	"github.com/switchupcb/copygen/cli/models"
)

// GENERATOR FUNCTION
// EDITABLE.
// DO NOT REMOVE.
// Function provides the generated code for each function.
func Function(function models.Function) string {
	return DefaultFunction(function)
}

// DefaultFunction provides generated code for a function using the default method.
func DefaultFunction(function models.Function) string {
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
		toComment += toType.Name + ", "
	}
	if len(toComment) != 0 {
		// remove last ", "
		toComment = toComment[:len(toComment)-2]
	}

	var fromComment string
	for _, fromType := range function.From {
		fromComment += fromType.Name + ", "
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
		parameters += toType.VariableName + " "
		if toType.Options.Pointer {
			parameters += "*"
		}
		if toType.Package != "" {
			parameters += toType.Package + "."
		}
		parameters += toType.Name + ", "
	}

	// Generate From-Type parameters
	for _, fromType := range function.From {
		parameters += fromType.VariableName + " "
		if fromType.Options.Pointer {
			parameters += "*"
		}
		if fromType.Package != "" {
			parameters += fromType.Package + "."
		}
		parameters += fromType.Name + ", "
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
		body += "// " + toType.Name + " fields\n"

		for _, toField := range toType.Fields {
			body += toType.VariableName + generateAssignment(toField)
			body += " = "
			fromField := toField.From
			if fromField.Convert != "" {
				body += fromField.Convert + "(" + fromField.Parent.VariableName + generateAssignment(fromField) + ")\n"
			} else {
				body += fromField.Parent.VariableName + generateAssignment(fromField) + "\n"
			}
		}
		body += "\n"
	}
	return body
}

// generateAssignment generates an assignment operation for the assignment of fields.
func generateAssignment(field *models.Field) string {
	if field.Of == nil {
		return "." + field.Name
	}
	return generateAssignment(field.Of) + "." + field.Name
}

// generateReturn generates a return statement for the function.
func generateReturn(function models.Function) string {
	return ""
}
