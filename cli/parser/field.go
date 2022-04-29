package parser

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/fatih/structtag"
	"github.com/switchupcb/copygen/cli/models"
)

// parseField parses a types.Type into a *models.Field recursively.
func parseField(typ types.Type) *models.Field {
	if cached, ok := fieldcache[typ.String()]; ok {
		return cached
	}

	field := new(models.Field)
	switch x := typ.(type) {

	// Named Types (Alias)
	// https://go.googlesource.com/example/+/HEAD/gotypes#named-types
	case *types.Named:
		// A named type is either:
		//   1. an alias (i.e `Placeholder` in `type Placeholder bool`)
		//   2. a struct (i.e `Account` in `type Account struct`)
		//   3. an interface (i.e `error` in `type error interface`)
		//   4. a collected type (i.e `domain.Account` in `[]domain.Account`)
		//
		// Underlying named types are only important in case 2,
		// when we need to parse extra information from the field.
		if xs, ok := x.Underlying().(*types.Struct); ok {

			// set the cache early to prevent issues with named cyclic structs.
			fieldcache[x.String()] = field
			structfield := parseField(xs)
			field.Fields = structfield.Fields
		}

		field.Definition = x.Obj().Name()
		setFieldImportAndPackage(field, x.Obj().Pkg())

	// Basic Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#basic-types
	case *types.Basic:
		field.Definition = x.Name()

	// Simple Composite Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#simple-composite-types
	case *types.Pointer:
		elemfield := parseField(x.Elem())

		// type aliases (including structs) must be deepcopied.
		if elemfield.IsAlias() {
			field = elemfield.Deepcopy(nil)
		}
		field.Definition = models.CollectionPointer + elemfield.Definition

	case *types.Array:
		field.Definition = "[" + fmt.Sprint(x.Len()) + "]" + parseField(x.Elem()).Definition

	case *types.Slice:
		field.Definition = models.CollectionSlice + parseField(x.Elem()).Definition

	case *types.Map:
		field.Definition = models.CollectionMap + "[" + parseField(x.Key()).Definition + "]" + parseField(x.Elem()).Definition

	case *types.Chan:
		field.Definition = models.CollectionChan + " " + parseField(x.Elem()).Definition

	// Function (without Receivers)
	// https://go.googlesource.com/example/+/HEAD/gotypes#function-and-method-types
	case *types.Signature:
		var definition strings.Builder

		// set the parameters.
		definition.WriteString(models.CollectionFunc + "(")
		for i := 0; i < x.Params().Len(); i++ {
			definition.WriteString(parseField(x.Params().At(i).Type()).Definition)
			if i+1 != x.Params().Len() {
				definition.WriteString(", ")
			}
		}
		definition.WriteString(")")

		// set the results.
		if x.Results().Len() >= 1 {
			definition.WriteString(" ")
		}
		if x.Results().Len() > 1 {
			definition.WriteString("(")
		}
		for i := 0; i < x.Results().Len(); i++ {
			definition.WriteString(parseField(x.Results().At(i).Type()).Definition)
			if i+1 != x.Results().Len() {
				definition.WriteString(", ")
			}
		}
		if x.Results().Len() > 1 {
			definition.WriteString(")")
		}

		field.Definition = definition.String()

	// Interface Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#interface-types
	case *types.Interface:
		if x.Empty() {
			field.Definition = x.String()
		} else {
			var definition strings.Builder
			definition.WriteString(models.CollectionInterface + "{")
			for i := 0; i < x.NumMethods(); i++ {
				definition.WriteString(parseField(x.Method(i).Type()).Definition + "; ")
			}
			definition.WriteString("}")
			field.Definition = definition.String()
		}

	// Struct Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#struct-types
	case *types.Struct:
		var definition strings.Builder
		definition.WriteString("struct{")
		for i := 0; i < x.NumFields(); i++ {
			// a deepcopy of subfield is returned and modified.
			subfield := parseField(x.Field(i).Type()).Deepcopy(nil)
			subfield.VariableName = "." + x.Field(i).Name()
			subfield.Name = x.Field(i).Name()
			setTags(subfield, x.Tag(i))
			subfield.Parent = field
			field.Fields = append(field.Fields, subfield)
			definition.WriteString(subfield.Name + " " + subfield.Definition + "; ")

			// all subfields are deepcopied with Fields[0].
			//
			// in order to correctly represent a deepcopied struct field,
			// we must point its fields back to the cached field.Fields,
			// which will eventually be filled.
			//
			// cachedsubfield.Fields are never modified.
			if cachedsubfield, ok := fieldcache[x.Field(i).String()]; ok {
				subfield.Fields = cachedsubfield.Fields
			}
		}
		definition.WriteString("}")
		field.Definition = definition.String()

	default:
		fmt.Printf("WARNING: could not parse type %v\n", x.String())
	}

	// do NOT cache collections.
	if !field.IsCollection() {
		fieldcache[typ.String()] = field
	}
	return field
}

// setFieldImportAndPackage sets the import and package of a field.
func setFieldImportAndPackage(field *models.Field, pkg *types.Package) {
	if pkg == nil {
		return
	}

	field.Import = pkg.Path()
	field.Package = pkg.Name()

	// field collections set collected types' packages in the field.Definition.
	// i.e map[*domain.Account]string
	if field.IsCollection() {
		field.Definition = field.Package + "." + field.Definition
		field.Import = ""
		field.Package = ""
	}
}

// setTags sets the tags for a field.
func setTags(field *models.Field, rawtag string) {
	// rawtag represents tags as they are defined (i.e `api:"id", json:"tag"`).
	tags, err := structtag.Parse(rawtag)
	if err != nil {
		fmt.Printf("WARNING: could not parse tag for field %v\n%v", field.FullName(), err)
	}

	field.Tags = make(map[string]map[string][]string, tags.Len())
	for _, tag := range tags.Tags() {
		field.Tags[tag.Key] = map[string][]string{
			tag.Name: tag.Options,
		}
	}
}
