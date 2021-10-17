package parser

import (
	"go/ast"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
)

// FieldSearch represents a search that uses Abstract Syntax Tree analysis to find the fields of a typefield.
type FieldSearch struct {
	// A key value cache used to prevent cyclic fields from unnecessary duplication or stack overflow.
	cache map[string]bool

	// The import of the field (used for the cache).
	Import string

	// The package name of the field.
	Package string

	// The name of the field.
	Name string

	// The actual typename of the field (i.e `DomainUser` in `User DomainUser`).
	Definition string

	// The parent of the field the FieldSearch will find.
	// In the context of the program, a top-level field with no parent is a TypeField.
	Parent *models.Field

	// The file that holds the type declaration for the field being searched.
	DecFile *ast.File

	// The options that pertain to a field (and its subfields).
	Options []Option

	// The current depth-level of the FieldSearch.
	Depth int

	// The maximum allowed depth-level of the FieldSearch.
	MaxDepth int
}

// SearchForField executes a field search by locating a field's type declaration, then its subfields.
func (p *Parser) SearchForField(fs *FieldSearch) (*models.Field, error) {
	cachename := fs.Import + fs.Package + fs.Name
	if cachedsearch, ok := p.fieldcache[cachename]; ok {
		if _, exists := fs.cache[cachename]; exists {
			// a cyclic field (with the same type as its parent) is never
			// shallow copied or assigned (unlike its parent or parent's fields).
			cyclicfield := &models.Field{
				VariableName: "." + cachedsearch.Name,
				Package:      cachedsearch.Package,
				Name:         cachedsearch.Name,
				Definition:   cachedsearch.Definition,
				Pointer:      cachedsearch.Pointer,
				Parent:       cachedsearch,
			}

			// depth is ignored for cyclic fields.
			setFieldOptions(cyclicfield, fs.Options)

			return cyclicfield, nil
		}

		fs.cache[cachename] = true

		return cachedsearch, nil
	}

	// setup the field
	field := &models.Field{
		Package:    fs.Package,
		Name:       fs.Name,
		Definition: fs.Definition,
	}

	if field.Definition == "" {
		// a TypeField definition is its name (i.e `Account` in `type Account struct`)
		field.Definition = field.Name
	}

	if fs.Parent != nil {
		field.VariableName = "." + fs.Name
		field.Parent = fs.Parent
	}

	setFieldOptions(field, fs.Options)
	fs.MaxDepth += field.Options.Depth

	// set the cache
	p.fieldcache[cachename] = field
	fs.cache[cachename] = true

	// find the fields of the main field if the max depth-level has not been reached.
	if fs.MaxDepth == 0 || fs.Depth < fs.MaxDepth {
		subfields, err := p.astSubfieldSearch(fs, field)
		if err != nil {
			return nil, err
		}

		field.Fields = subfields
	}

	return field, nil
}

// astSubfieldSearch searches through an ast.Typespec for sub-fields.
func (p *Parser) astSubfieldSearch(fs *FieldSearch, typefield *models.Field) ([]*models.Field, error) {
	// the original setup file (i.e setup.go) is used to locate the file location of a field's actual type declaration.
	td, err := p.astLocateTypeDecl(&Locater{
		SetupFile:  fs.DecFile,
		Package:    typefield.Package,
		Definition: typefield.Definition,
	})
	if err != nil {
		return nil, err
	}

	// then that data is parsed into subfield information.
	var subfields []*models.Field

	switch x := td.Package.TypesInfo.Types[td.TypeSpec.Type].Type.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xField := x.Field(i)

			// create a new typefield if a subfield is a custom type
			if !isBasic(xField.Type()) {
				parsedDefinition := p.parseDefinition(xField.Type().String())
				if parsedDefinition.err != nil {
					return nil, parsedDefinition.err
				}

				// Search for the subfields of the subfield
				subfield, err := p.SearchForField(&FieldSearch{
					Import:     parsedDefinition.imprt,
					Name:       xField.Name(),
					Package:    parsedDefinition.pkg,
					Definition: parsedDefinition.typename,
					Parent:     typefield,
					DecFile:    td.File,
					Options:    fs.Options,
					Depth:      fs.Depth + 1,
					MaxDepth:   fs.MaxDepth,
					cache:      fs.cache,
				})
				if err != nil {
					return nil, err
				}

				subfields = append(subfields, subfield)
			} else {
				subfield := &models.Field{
					VariableName: "." + xField.Name(),
					Name:         xField.Name(),
					Definition:   xField.Type().String(),
					Parent:       typefield,
				}
				setFieldOptions(subfield, fs.Options)
				subfields = append(subfields, subfield)
			}
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
				Parent:       typefield,
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
		return true
	case *types.Map:
		return true
	case *types.Pointer:
		return isBasic(x.Elem())
	default:
		return false
	}
}

// setFieldOptions sets a field's (and its subfields) options.
func setFieldOptions(field *models.Field, options []Option) {
	setConvertOption(field, options)
	setDeepcopyOption(field, options)
	setDepthOption(field, options)
	setMapOption(field, options)
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
