package parser

import (
	"fmt"
	"go/ast"
	"strconv"

	"github.com/switchupcb/copygen/cli/models"
)

type parsedTypes struct {
	fromTypes []models.Type
	toTypes   []models.Type
}

// parseTypes parses an ast.Field (of type func) for to-types and from-types.
func (p *Parser) parseTypes(function *ast.Field, fieldsearcher *FieldSearcher) (parsedTypes, error) {
	var result parsedTypes

	fn, ok := function.Type.(*ast.FuncType)
	if !ok {
		return result, fmt.Errorf("an error occurred parsing the types of function %v at Line %d", parseMethodForName(function), p.Fileset.Position(function.Pos()).Line)
	}

	fromTypes, err := p.parseFieldList(fn.Params.List, fieldsearcher) // (incoming) parameters "non-nil"
	if err != nil {
		return result, err
	}
	var toTypes []models.Type
	if fn.Results != nil {
		toTypes, err = p.parseFieldList(fn.Results.List, fieldsearcher) // (outgoing) results "or nil"
		if err != nil {
			return result, err
		}
	}
	if len(fromTypes) == 0 {
		return result, fmt.Errorf("function %v at Line %d has no types to copy from", parseMethodForName(function), p.Fileset.Position(function.Pos()).Line)
	} else if len(toTypes) == 0 {
		return result, fmt.Errorf("function %v at Line %d has no types to copy to", parseMethodForName(function), p.Fileset.Position(function.Pos()).Line)
	}

	// assign variable names and determine the definition and sub-fields of each type
	paramMap := make(map[string]bool)
	for i := 0; i < len(fromTypes); i++ {
		fromTypes[i].Field.VariableName = createVariable(paramMap, "f"+fromTypes[i].Field.Name, 0)
	}
	for i := 0; i < len(toTypes); i++ {
		toTypes[i].Field.VariableName = createVariable(paramMap, "t"+toTypes[i].Field.Name, 0)
	}

	result.fromTypes = fromTypes
	result.toTypes = toTypes

	return result, nil
}

// parseFieldList parses an Abstract Syntax Tree field list for a type's fields.
func (p *Parser) parseFieldList(fieldlist []*ast.Field, fieldsearcher *FieldSearcher) ([]models.Type, error) {
	types := make([]models.Type, 0, len(fieldlist))

	for _, astfield := range fieldlist {
		field, err := p.parseTypeField(astfield, fieldsearcher)
		if err != nil {
			return nil, err
		}
		types = append(types, models.Type{Field: field})
	}
	return types, nil
}

// parseTypeField parses a function *ast.Field into a field model.
func (p *Parser) parseTypeField(field *ast.Field, fieldsearcher *FieldSearcher) (*models.Field, error) {
	parsed := parseASTFieldName(field)
	if parsed.name == "" {
		return nil, fmt.Errorf("unexpected field expression %v in the Abstract Syntax Tree", field)
	}

	mField, err := fieldsearcher.SearchForTypeField(p.SetupFile, p.Imports[parsed.pkg], parsed.pkg, parsed.name)
	if err != nil {
		return nil, err
	}
	mField.Pointer = parsed.ptr

	return mField, nil
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
