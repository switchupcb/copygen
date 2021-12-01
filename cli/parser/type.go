package parser

import (
	"fmt"
	"go/types"
	"path"
	"strconv"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

type parsedTypes struct {
	fromTypes []models.Type
	toTypes   []models.Type
}

func (p *Parser) packageNameByImport(imp string) string {
	if p.gen.ImportsByPath[imp] == "" {
		k := 0
		bn := path.Base(imp)
		// Fix for packages like `yaml.v3` with dot in name. Only `yaml` should be taken.
		if strings.Index(bn, ".") > 0 {
			bn = bn[0:strings.Index(bn, ".")]
		}
		for k = 0; p.gen.ImportsByName[bn+strconv.Itoa(k)] != ""; k++ {
		}
		name := bn
		if k > 0 {
			name += strconv.Itoa(k)
		}
		p.gen.ImportsByName[name] = imp
		p.gen.ImportsByPath[imp] = name
	}
	return p.gen.ImportsByPath[imp]
}

// parseTypes parses an ast.Field (of type func) for to-types and from-types.
func (p *Parser) parseTypes(function *types.Signature, options []Option) (parsedTypes, error) {
	var result parsedTypes

	fromTypes, err := p.parseFieldList(function.Params(), options) // (incoming) parameters "non-nil"
	if err != nil {
		return result, err
	}

	var toTypes []models.Type

	if function.Results() != nil {
		toTypes, err = p.parseFieldList(function.Results(), options) // (outgoing) results "or nil"
		if err != nil {
			return result, err
		}
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
func (p *Parser) parseFieldList(fieldlist *types.Tuple, options []Option) ([]models.Type, error) {
	types := make([]models.Type, 0, fieldlist.Len())

	for i := 0; i < fieldlist.Len(); i++ {
		field, err := p.parseTypeField(fieldlist.At(i), options)
		field.Package = p.packageNameByImport(field.Import)
		if err != nil {
			return nil, err
		}

		types = append(types, models.Type{Field: field})
	}

	return types, nil
}

func (p *Parser) unwrapPointer(t types.Type) (string, *types.Named) {
	out := ""
	if subType, ok := t.(*types.Pointer); ok {
		out = "*"
		resp, newT := p.unwrapPointer(subType.Elem())
		out += resp
		t = newT
	}
	return out, t.(*types.Named)
}

// parseTypeField parses a function *ast.Field into a field model.
func (p *Parser) parseTypeField(field *types.Var, options []Option) (*models.Field, error) {
	strPtr, ptr := p.unwrapPointer(field.Type())
	parsed := parsedFieldName{
		pkg:  field.Pkg().Name(),
		name: ptr.Obj().Name(),
		ptr:  strPtr,
	}
	if parsed.name == "" {
		return nil, fmt.Errorf("unexpected field expression %v in the Abstract Syntax Tree", field)
	}

	typefield, err := p.SearchForField(&FieldSearch{
		Options: options,
		cache:   make(map[string]bool),
	}, field.Type())
	if err != nil {
		return nil, fmt.Errorf("an error occurred while searching for the top-level Field %q of package %q.\n%v", parsed.name, parsed.pkg, err)
	}

	typefield.Pointer = parsed.ptr

	return typefield, nil
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
