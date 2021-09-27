// package template implements an interpreter used to provide customizable functions to the generator.
package template

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/models"
)

// Function determines the func to generate function code.
func Function(gen *models.Generator) (string, error) {
	var functions string

	// determine the method to analyze each function.
	if gen.Template.Funcpath == "" {
		for _, function := range gen.Functions {
			functions += defaultFunction(&function) + "\n"
		}
		return functions, nil
	}
	return "", fmt.Errorf("Templates are temporarily unsupported.")

	fn, err := interpretFunction(gen)
	if err != nil {
		return "", err
	}
	for _, function := range gen.Functions {
		functions += fn(&function) + "\n"
	}
	return functions, nil
}

// defaultFunction creates the header of the generated file using the default method.
func defaultFunction(function *models.Function) string {
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
func generateComment(function *models.Function) string {
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
func generateSignature(function *models.Function) string {
	sig := "func " + function.Name + "(" + generateParameters(function) + ") {"
	return sig
}

// generateParameters generates the parameters of a function.
func generateParameters(function *models.Function) string {
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
func generateBody(function *models.Function) string {
	var body string

	// Assign fields to ToType(s)
	for _, toType := range function.To {
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

func generateReturn(function *models.Function) string {
	return ""
}

// interpretFunction creates the header of the generated file using an interpreted template file.
func interpretFunction(gen *models.Generator) (func(f *models.Function) string, error) {
	fn, err := interpreter.InterpretFunc(gen.Loadpath, gen.Template.Funcpath, "generator.Function")
	if err != nil {
		return nil, err
	}

	// run the interpreted function.
	return func(function *models.Function) string {
		return fn(function)
	}, nil
}
