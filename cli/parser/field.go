package parser

import (
	"go/ast"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
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
	Options []*options.Option

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
			p.setFieldOptions(cyclicfield, fs.Options)

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

	p.setFieldOptions(field, fs.Options)
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
				p.setFieldOptions(subfield, fs.Options)
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
			p.setFieldOptions(subfield, fs.Options)
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
func (p *Parser) setFieldOptions(field *models.Field, opts []*options.Option) {
	options.SetConvert(field, p.ConvertOptions)
	options.SetDeepcopy(field, opts)
	options.SetDepth(field, opts)
	options.SetMap(field, opts)
}
