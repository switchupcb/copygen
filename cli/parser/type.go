package parser

import (
	"fmt"
	"go/types"
	"strconv"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

type parsedTypes struct {
	fromTypes []models.Type
	toTypes   []models.Type
}

// parseTypes parses a types.Func's parameters for from-types and results for to-types.
func parseTypes(method *types.Func, options []*options.Option) (parsedTypes, error) {
	var result parsedTypes

	signature, ok := method.Type().(*types.Signature)
	if !ok {
		return result, fmt.Errorf("impossible")
	}

	if signature.Params().Len() == 0 {
		return result, fmt.Errorf("function %v has no types to copy from", method.Name())
	} else if signature.Results().Len() == 0 {
		return result, fmt.Errorf("function %v has no types to copy to", method.Name())
	}

	var err error
	result.fromTypes, err = parseTypeField(signature.Params(), options)
	if err != nil {
		return result, fmt.Errorf("an error occurred while parsing a from type parameter in %v\n%w", method.Name(), err)
	}

	result.toTypes, err = parseTypeField(signature.Results(), options)
	if err != nil {
		return result, fmt.Errorf("an error occurred while parsing a from type parameter in %v\n%w", method.Name(), err)
	}

	setVariableNames(result.fromTypes, "f")
	setVariableNames(result.toTypes, "t")

	return result, nil
}

// parseTypeField parses a *types.Tuple into a *models.Type (that points to a *models.Field).
func parseTypeField(vars *types.Tuple, fieldoptions []*options.Option) ([]models.Type, error) {
	types := make([]models.Type, vars.Len())
	for i := 0; i < vars.Len(); i++ {

		// create a top-level field (fieldParser parent = nil).
		fp := fieldParser{options: fieldoptions, cyclic: make(map[string]bool)}
		field := fp.parseField(vars.At(i).Type())
		if field == nil {
			return nil, fmt.Errorf("an error occurred parsing a type field parameter %v", vars.At(i).String())
		}

		field.Name = vars.At(i).Name()
		types[i] = models.Type{
			Field: field,
		}
	}

	return types, nil
}

// setVariableNames sets the variable names for a list of type fields.
func setVariableNames(types []models.Type, precedent string) {
	paramMap := make(map[string]bool)
	for i := 0; i < len(types); i++ {
		types[i].Field.VariableName = createVariable(paramMap, precedent+types[i].Field.VariableName[1:], 0)
	}
}

// createVariable generates a valid variable name for a 'set' of parameters.
func createVariable(parameters map[string]bool, typename string, occurrence int) string {
	if occurrence < 0 {
		createVariable(parameters, typename, 0)
	}

	varName := typename[:2]
	if occurrence > 0 {
		varName += strconv.Itoa(occurrence + 1)
	}

	if _, exists := parameters[varName]; exists {
		createVariable(parameters, typename, occurrence+1)
	}

	return varName
}
