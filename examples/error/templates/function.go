package generator

import (
	"github.com/switchupcb/copygen/cli/models"
)

// Function generates code for a function.
func Function(f *models.Function) string {
	var function string

	// comment
	function += generateComment(f) + "\n"

	// signature
	function += generateSignature(f) + "\n"

	// body
	function += generateBody(f) + "\n"

	// return
	function += generateReturn(f) + "\n"

	// end of function
	function += "}"
	return function
}

// generateComment generates a function comment.
func generateComment(f *models.Function) string {
	var tComment string
	for _, toType := range f.To {
		tComment += toType.Name + ", "
	}
	if len(tComment) != 0 {
		// remove last ", "
		tComment = tComment[:len(tComment)-2]
	}

	var fComment string
	for _, fromType := range f.From {
		fComment += fromType.Name + ", "
	}
	if len(fComment) != 0 {
		// remove last ", "
		fComment = fComment[:len(fComment)-2]
	}

	return "// " + f.Name + " copies a " + fComment + " to a " + tComment + "."
}

// generateSignature generates a function's signature.
func generateSignature(f *models.Function) string {
	s := "func " + f.Name + "(" + generateParameters(f) + ") {"
	return s
}

// generateParameters generates the parameters of a function.
func generateParameters(f *models.Function) string {
	var parameters string

	// Generate To-Type parameters
	for _, toType := range f.To {
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
	for _, fromType := range f.From {
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
func generateBody(f *models.Function) string {
	var body string

	// Assign fields to ToType(s)
	for _, toType := range f.To {
		body += "// " + toType.Name + " fields\n"

		for _, toField := range toType.Fields {
			// toField
			body += toType.VariableName + "." + toField.Name + " = "

			// fromField
			if toField.Convert != "" {
				body += toField.Convert + "(" + toField.Parent.VariableName + "." + toField.From.Name + ")\n"
			} else {
				body += toField.From.Parent.VariableName + "." + toField.From.Name + "\n"
			}
		}
		body += "\n"
	}
	return body
}

func generateReturn(f *models.Function) string {
	return ""
}
