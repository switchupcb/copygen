package parser

import (
	"fmt"
	"go/types"

	"github.com/fatih/structtag"
	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

// fieldParser represents the parameters required to parse a types.Type into a *models.Field.
type fieldParser struct {
	// field represents the current field being built.
	field *models.Field

	// parent represents the parent of the field parse.
	parent *models.Field

	// cyclic is a key value cache used to prevent cyclic fields from unnecessary duplication or stack overflow.
	cyclic map[string]bool

	// container represents a field's container.
	container string

	// options represents the field options defined above the models.Function
	options []*options.Option
}

// parseField parses a types.Type into a *models.Field.
func (fp fieldParser) parseField(typ types.Type) *models.Field {
	if fp.field == nil {
		fp.field = &models.Field{Parent: fp.parent}
	}

	switch x := typ.(type) {

	// Basic Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#basic-types
	case *types.Basic:
		setFieldVariableName(fp.field, "."+x.Name())
		fp.field.Definition = x.Name()

	// Named Types (Alias)
	// https://go.googlesource.com/example/+/HEAD/gotypes#named-types
	case *types.Named:
		setFieldImportAndPackage(fp.field, x.Obj().Pkg().Path(), x.Obj().Pkg().Name())
		setFieldVariableName(fp.field, "."+x.Obj().Name())
		setFieldName(fp.field, x.Obj().Name())
		return fp.parseField(x.Underlying())

	// Simple Composite Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#simple-composite-types
	case *types.Pointer:
		fp.field.Container += "*"
		return fp.parseField(x.Elem())

	case *types.Array:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "[" + fmt.Sprint(x.Len()) + "]"

	case *types.Slice:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "[]"

	case *types.Map:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "map"

	case *types.Chan:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "chan"

	// Struct Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#struct-types
	case *types.Struct:
		fp.field.Collection = "struct"
		for i := 0; i < x.NumFields(); i++ {
			subfield := &models.Field{
				VariableName: "." + x.Field(i).Name(),
				Name:         x.Field(i).Name(),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Field(i).Pkg().Path(), x.Field(i).Pkg().Name())
			setTags(subfield, x.Tag(i))

			// a cyclic subfield (with the same type as its parent) is never fully assigned.
			if !fp.cyclic[x.Field(i).String()] {
				subfieldParser := &fieldParser{
					field:     subfield,
					parent:    nil,
					container: "",
					options:   fp.options,
					cyclic:    fp.cyclic,
				}

				// sets the definition, container, and fields.
				fp.cyclic[x.Field(i).String()] = true
				subfield = subfieldParser.parseField(x.Field(i).Type())
			}

			fp.field.Fields = append(fp.field.Fields, subfield)
		}

	// Function
	// https://go.googlesource.com/example/+/HEAD/gotypes#function-and-method-types
	case *types.Signature:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()

	// Interface Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#interface-types
	case *types.Interface:
		fp.field.Collection = "interface"
		for i := 0; i < x.NumMethods(); i++ {
			subfield := &models.Field{
				VariableName: "." + x.Method(i).Name(),
				Name:         x.Method(i).Name(),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Method(i).Pkg().Path(), x.Method(i).Pkg().Name())

			subfieldParser := &fieldParser{
				field:     subfield,
				parent:    nil,
				container: "",
				options:   fp.options,
				cyclic:    fp.cyclic,
			}

			// sets the definition, container, and fields.
			subfield = subfieldParser.parseField(x.Method(i).Type())
			fp.field.Fields = append(fp.field.Fields, subfield)
		}
	}

	options.SetFieldOptions(fp.field, fp.options)
	filterFieldDepth(fp.field, fp.field.Options.Depth, 0)
	fp.cyclic[typ.String()] = true
	return fp.field
}

// setFieldVariableName sets a field's variable name.
func setFieldVariableName(field *models.Field, varname string) {
	if field.VariableName == "" {
		field.VariableName = varname
	}
}

// setFieldName sets a field's name.
func setFieldName(field *models.Field, name string) {
	if field.Name == "" {
		field.Name = name
	} else {
		field.Definition = name
	}
}

// setFieldImportAndPackage sets the import and package of a field.
func setFieldImportAndPackage(field *models.Field, path string, varname string) {
	field.Import = path
	if ignorepkgpath != path {
		if _, ok := aliasImportMap[path]; ok {
			field.Package = aliasImportMap[path]
		} else {
			field.Package = varname
		}
	}
}

// setTags sets the tags for a field.
func setTags(field *models.Field, rawtag string) {
	// rawtag represents tags as they are defined (i.e `api:"id", json:"tag"`).
	tags, err := structtag.Parse(rawtag)
	if err != nil {
		fmt.Printf("WARNING: could not parse tag for field %v\n%v", field.FullName(""), err)
	}

	if field.Tags == nil {
		field.Tags = make(map[string]map[string][]string, tags.Len())
	}

	for _, tag := range tags.Tags() {
		field.Tags[tag.Key] = map[string][]string{
			tag.Name: tag.Options,
		}
	}
}

// filterFieldDepth filters a field's fields according to it's depth level.
func filterFieldDepth(field *models.Field, maxdepth, curdepth int) {
	if maxdepth == 0 {
		return
	}

	if maxdepth < 0 || maxdepth <= curdepth {
		field.Fields = make([]*models.Field, 0)
		return
	}

	for _, f := range field.Fields {
		filterFieldDepth(f, maxdepth+f.Options.Depth, curdepth+1)
	}
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
