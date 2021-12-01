package parser

import (
	"go/types"
	"reflect"

	"github.com/switchupcb/copygen/cli/models"
)

// FieldSearch represents a search that uses Abstract Syntax Tree analysis to find the fields of a typefield.
type FieldSearch struct {
	// A key value cache used to prevent cyclic fields from unnecessary duplication or stack overflow.
	cache map[string]bool

	// The options that pertain to a field (and its subfields).
	Options []Option

	// The current depth-level of the FieldSearch.
	Depth int

	// The maximum allowed depth-level of the FieldSearch.
	MaxDepth int
}

// SearchForField executes a field search by locating a field's type declaration, then its subfields.
func (p *Parser) SearchForField(fs *FieldSearch, typ types.Type) (*models.Field, error) {
	parsedDefinition := p.parseDefinition(typ.String())
	// We must reset cache otherwise second function with same methods will fail its execution .
	p.fieldcache = map[string]*models.Field{}
	// setup the field
	field := &models.Field{
		Name:           parsedDefinition.typename,
		VariableName:   parsedDefinition.typename,
		Import:         parsedDefinition.imprt,
		Definition:     parsedDefinition.typename,
		OrigDefinition: parsedDefinition.typename,
		Package:        p.packageNameByImport(parsedDefinition.imprt),
	}

	setFieldOptions(field, fs.Options)
	fs.MaxDepth += field.Options.Depth

	// find the fields of the main field if the max depth-level has not been reached.
	if fs.MaxDepth == 0 || fs.Depth < fs.MaxDepth {
		_, typ = p.unwrapPointer(typ)
		subfields, err := p.subfieldSearch(typ.(*types.Named).Obj().Type().Underlying(), fs, field)
		if err != nil {
			return nil, err
		}

		field.Fields = subfields
	}

	return field, nil
}

// subfieldSearch searches through an types.Type for sub-fields.
func (p *Parser) subfieldSearch(td types.Type, fs *FieldSearch, parent *models.Field) ([]*models.Field, error) {
	// then that data is parsed into subfield information.
	subfields := make([]*models.Field, 0)

	switch x := td.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xField := x.Field(i)
			xTag := x.Tag(i)
			cachename := x.String() + "." + xField.Name()
			if cachedsearch, ok := p.fieldcache[cachename]; ok {
				if _, exists := fs.cache[cachename]; exists {
					// a cyclic field (with the same type as its parent) is never
					// shallow copied or assigned (unlike its parent or parent's fields).
					cyclicfield := &(*cachedsearch)
					// depth is ignored for cyclic fields.
					setFieldOptions(cyclicfield, fs.Options)
					subfields = append(subfields, cyclicfield)
				}
				fs.cache[cachename] = true
				continue
			}
			// create a new typefield if a subfield is a custom type
			parsedDefinition := p.parseDefinition(xField.Type().String())
			tp := xField.Type()
			def := xField.Type().String()
			origDef := xField.Type().Underlying().String()
			if sl, ok := xField.Type().(*types.Slice); ok {
				origDef = sl.Elem().Underlying().String()
				parsedDefinition := p.parseDefinition(sl.Elem().String())
				def = parsedDefinition.typename
				tp = sl.Elem()
			}
			subfield := &models.Field{
				VariableName:   "." + xField.Name(),
				Name:           xField.Name(),
				Definition:     def,
				OrigDefinition: origDef,
				Import:         parsedDefinition.imprt,
				Parent:         parent,
				Tags:           reflect.StructTag(xTag),
				ContainerType:  parsedDefinition.containerType,
				Pointer:        parsedDefinition.pointer,
			}
			if parent == nil {
				subfield.VariableName = "." + xField.Name()
			}
			if parsedDefinition.imprt != "" {
				subfield.Package = p.packageNameByImport(parsedDefinition.imprt)
				subfield.Definition = parsedDefinition.typename
			}

			if !isBasic(tp) {
				if tpn, ok := tp.(*types.Pointer); ok {
					tp = tpn.Elem()
				}
				subfield.OrigDefinition = parsedDefinition.typename
				if parsedDefinition.err != nil {
					return nil, parsedDefinition.err
				}
				// Search for the subfields of the subfield
				subfield.OrigDefinition = origDef
				if _, ok := xField.Type().Underlying().(*types.Struct); ok {
					var err error
					subfield.Fields, err = p.subfieldSearch(xField.Type().Underlying(), fs, subfield)
					if err != nil {
						return nil, err
					}
				}
				subfields = append(subfields, subfield)
				subfields = append(subfields, subfield.Fields...)
			} else {
				subfields = append(subfields, subfield)
			}
			setFieldOptions(subfield, fs.Options)
			// set the cache
			p.fieldcache[cachename] = subfield
			fs.cache[cachename] = true
		}
	// interfaces have method fields
	case *types.Interface:
		for i := 0; i < x.NumMethods(); i++ {
			xMethod := x.Method(i)
			// interface functions are declared in the same scope as the interface
			subfield := &models.Field{
				VariableName: "." + xMethod.Name() + "(%)",
				Package:      xMethod.Pkg().Name(),
				Name:         xMethod.Name(),
				Definition:   xMethod.Type().String(),
				Parent:       parent,
			}
			setFieldOptions(subfield, fs.Options)
			subfields = append(subfields, subfield)
		}
	default: // if no fields are present, this is a basic type.
	}

	return subfields, nil
}

// isBasic determines whether a type is a basic (non-custom) type.
func isBasic(t types.Type) bool {
	switch x := t.(type) {
	case *types.Basic:
		return true
	case *types.Slice:
		return isBasic(x.Elem())
	case *types.Map:
		return true
	case *types.Pointer:
		return isBasic(x.Elem())
	default:
		if x == x.Underlying() {
			return false
		}
		return isBasic(x.Underlying())
	}
}

// setFieldOptions sets a field's (and its subfields) options.
func setFieldOptions(field *models.Field, options []Option) {
	setConvertOption(field, options)
	setDeepcopyOption(field, options)
	setDepthOption(field, options)
	setMapOption(field, options)
	setTagOption(field, options)
}

// setConvertOption sets a field's convert option.
func setConvertOption(field *models.Field, options []Option) {
	// A convert option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == categoryConvert && options[i].Regex[1].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Convert = value
				break
			}
		}
	}
}

// setDeepcopyOption sets a field's deepcopy option.
func setDeepcopyOption(field *models.Field, options []Option) {
	// A deepcopy option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == categoryDeepCopy && options[i].Regex[0].MatchString(field.FullName("")) {
			field.Options.Deepcopy = true
			break
		}
	}
}

// setDepthOption sets a field's depth option.
func setDepthOption(field *models.Field, options []Option) {
	// A depth option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == categoryDepth && options[i].Regex[0].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(int); ok {
				// Automatch all is on by default; if a user specifies 0 depth-level, guarantee it.
				if value == 0 {
					value = -1
				}

				field.Options.Depth = value

				break
			}
		}
	}
}

// setMapOption sets a field's deepcopy option.
func setMapOption(field *models.Field, options []Option) {
	// A map option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == categoryMap && options[i].Regex[0].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Map = value
				break
			}
		}
	}
}

func setTagOption(field *models.Field, options []Option) {
	// A tag option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == categoryCommonTag && options[i].Regex[0].MatchString(field.Package+"."+field.Name) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Tag = value
				break
			}
		}
	}
}
