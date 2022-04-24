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
	cyclic map[string]*models.Field

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
		setDefinition(fp.field, x.Name())

	// Named Types (Alias)
	// https://go.googlesource.com/example/+/HEAD/gotypes#named-types
	case *types.Named:
		setFieldImportAndPackage(fp.field, x.Obj().Pkg())
		setFieldVariableName(fp.field, "."+x.Obj().Name())
		setDefinition(fp.field, x.Obj().Name())

		// do NOT parse named types in a collection.
		if !fp.field.IsCollection() {
			return fp.parseField(x.Underlying())
		}

	// Simple Composite Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#simple-composite-types
	case *types.Pointer:
		if fp.field.Definition == "" && fp.field.Pointer == "" {
			fp.field.Pointer = models.Pointer
		} else {
			setDefinition(fp.field, models.CollectionPointer)
		}
		return fp.parseField(x.Elem())

	case *types.Array:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		setDefinition(fp.field, "["+fmt.Sprint(x.Len())+"]")
		return fp.parseField(x.Elem())

	case *types.Slice:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		setDefinition(fp.field, models.CollectionSlice)
		return fp.parseField(x.Elem())

	case *types.Map:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		setDefinition(fp.field, models.CollectionMap+"[")
		_ = fp.parseField(x.Key())
		setDefinition(fp.field, "]")
		return fp.parseField(x.Elem())

	case *types.Chan:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		setDefinition(fp.field, models.CollectionChan+" ")
		return fp.parseField(x.Elem())

	// Function (without Receivers)
	// https://go.googlesource.com/example/+/HEAD/gotypes#function-and-method-types
	case *types.Signature:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))

		// set the parameters.
		setDefinition(fp.field, models.CollectionFunc+"(")
		for i := 0; i < x.Params().Len(); i++ {
			_ = fp.parseField(x.Params().At(i).Type())
			if i+1 != x.Params().Len() {
				setDefinition(fp.field, ", ")
			}
		}
		setDefinition(fp.field, ")")

		// set the results.
		if x.Results().Len() > 1 {
			setDefinition(fp.field, " (")
		}
		for i := 0; i < x.Results().Len(); i++ {
			_ = fp.parseField(x.Results().At(i).Type())
			if i+1 != x.Results().Len() {
				setDefinition(fp.field, ", ")
			}
		}
		if x.Results().Len() > 1 {
			setDefinition(fp.field, ")")
		}

	// Struct Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#struct-types
	case *types.Struct:
		fp.field.Container = models.ContainerStruct
		for i := 0; i < x.NumFields(); i++ {
			if subfield, ok := fp.cyclic[x.Field(i).String()]; ok {
				fp.field.Fields = append(fp.field.Fields, subfield)
				continue
			}

			// parse a new field.
			subfield := &models.Field{
				VariableName: "." + x.Field(i).Name(),
				Name:         x.Field(i).Name(),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Field(i).Pkg())
			setTags(subfield, x.Tag(i))
			subfieldParser := &fieldParser{
				field:     subfield,
				parent:    nil,
				container: "",
				options:   fp.options,
				cyclic:    fp.cyclic,
			}

			// sets the definition, container, and fields.
			fp.cyclic[x.Field(i).String()] = subfield
			subfield = subfieldParser.parseField(x.Field(i).Type())
			fp.field.Fields = append(fp.field.Fields, subfield)
		}

	// Interface Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#interface-types
	case *types.Interface:
		fp.field.Container = models.ContainerInterface
		for i := 0; i < x.NumMethods(); i++ {
			if subfield, ok := fp.cyclic[x.Method(i).String()]; ok {
				fp.field.Fields = append(fp.field.Fields, subfield)
				continue
			}

			// parse a new field.
			subfield := &models.Field{
				VariableName: "." + x.Method(i).Name(),
				Name:         x.Method(i).Name(),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Method(i).Pkg())
			subfieldParser := &fieldParser{
				field:     subfield,
				parent:    nil,
				container: "",
				options:   fp.options,
				cyclic:    fp.cyclic,
			}

			// sets the definition, container, and fields.
			fp.cyclic[x.Method(i).String()] = subfield
			subfield = subfieldParser.parseField(x.Method(i).Type())
			fp.field.Fields = append(fp.field.Fields, subfield)
		}
	}

	options.SetFieldOptions(fp.field, fp.options)
	filterFieldDepth(fp.field, fp.field.Options.Depth, 0)
	fp.cyclic[typ.String()] = fp.field
	return fp.field
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

// setFieldImportAndPackage sets the import and package of a field.
func setFieldImportAndPackage(field *models.Field, pkg *types.Package) {
	if pkg == nil {
		return
	}

	field.Import = pkg.Path()
	if ignorepkgpath != field.Import {
		if _, ok := aliasImportMap[field.Import]; ok {
			field.Package = aliasImportMap[field.Import]
		} else {
			field.Package = pkg.Name()
		}
	}

	if field.IsCollection() {
		setDefinition(field, field.Package+".")
		field.Package = ""
	}
}

// setFieldVariableName sets a field's variable name.
func setFieldVariableName(field *models.Field, varname string) {
	if field.VariableName == "" {
		field.VariableName = varname
	}
}

// setDefinition sets a field's definition.
func setDefinition(field *models.Field, def string) {
	switch {
	case field.Definition == "":
		field.Definition = def
	case field.IsCollection():
		field.Definition += def
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
