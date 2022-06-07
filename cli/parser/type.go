package parser

import (
	"fmt"
	"go/types"
	"strconv"

	"github.com/switchupcb/copygen/cli/models"
)

type parsedTypes struct {
	fromTypes []models.Type
	toTypes   []models.Type
}

// parseTypes parses a types.Func's parameters for from-types and results for to-types.
func parseTypes(method *types.Func) (parsedTypes, error) {
	var result parsedTypes

	signature, ok := method.Type().(*types.Signature)
	if !ok {
		return result, fmt.Errorf("impossible")
	}

	result.fromTypes = parseTypeField(signature.Params())
	setVariableNames(result.fromTypes, "f")

	result.toTypes = parseTypeField(signature.Results())
	setVariableNames(result.toTypes, "t")

	return result, nil
}

// parseTypeField parses a *types.Tuple into a *models.Type (that points to a *models.Field).
func parseTypeField(vars *types.Tuple) []models.Type {
	types := make([]models.Type, vars.Len())
	for i := 0; i < vars.Len(); i++ {
		field := parseField(vars.At(i).Type()).Deepcopy(nil)
		field.Name = vars.At(i).Name()
		if !field.IsPointer() {
			field.VariableName = "." + alphastring(field.Definition)
		}
		types[i] = models.Type{Field: field}
	}

	return types
}

// setVariableNames sets the variable names for a list of type fields.
func setVariableNames(types []models.Type, precedent string) {
	paramMap := make(map[string]bool)
	for i := 0; i < len(types); i++ {
		types[i].Field.VariableName = createVariable(paramMap, precedent+types[i].Field.VariableName[1:], 0)
		paramMap[types[i].Field.VariableName] = true
	}
}

// createVariable generates a valid variable name for a 'set' of parameters.
func createVariable(parameters map[string]bool, name string, occurrence int) string {
	if occurrence < 0 {
		createVariable(parameters, name, 0)
	}

	// assume a precedent (i.e `t`) and variable (min size = 1) has been passed.
	varname := name[:2]
	if occurrence > 0 {
		varname += strconv.Itoa(occurrence)
	}

	if parameters[varname] {
		return createVariable(parameters, name, occurrence+1)
	}

	return varname
}

// alphastring only returns alphabetic characters (English) in a string.
func alphastring(s string) string {
	bytes := []byte(s)
	i := 0
	for _, b := range bytes {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == ' ' {
			bytes[i] = b
			i++
		}
	}

	return string(bytes[:i])
}
