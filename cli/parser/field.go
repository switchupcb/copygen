package parser

import (
	"fmt"
	"go/types"
	"strconv"
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
		// set the cache early to prevent issues with named cyclic types.
		fieldcache[x.String()] = field

		// A named type is either:
		//   1. an alias (i.e `Placeholder` in `type Placeholder bool`)
		//   2. a struct (i.e `Account` in `type Account struct`)
		//   3. an interface (i.e `error` in `type error interface`)
		//   4. a collected type (i.e `domain.Account` in `[]domain.Account`)
		//
		// Underlying named types are only important in case 2,
		// when we need to parse extra information from the field.
		if xs, ok := x.Underlying().(*types.Struct); ok {
			structfield := parseField(xs)
			field.Fields = structfield.Fields
		} else {
			field.Underlying = parseField(x.Underlying())
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

		// type aliases (including structs) must be deepcopied
		// in order to match underlying fields.
		if elemfield.IsAlias() {
			deepfield := elemfield.Deepcopy(nil)
			field.Fields = deepfield.Fields
		}

		field.Definition = models.CollectionPointer + collectedDefinition(elemfield)
		field.VariableName = "." + alphastring(elemfield.Definition)

	case *types.Array:
		field.Definition = "[" + strconv.FormatInt(x.Len(), 10) + "]" + collectedDefinition(parseField(x.Elem()))

	case *types.Slice:
		field.Definition = models.CollectionSlice + collectedDefinition(parseField(x.Elem()))

	case *types.Map:
		field.Definition = models.CollectionMap + "[" + collectedDefinition(parseField(x.Key())) + "]" + collectedDefinition(parseField(x.Elem()))

	case *types.Chan:
		field.Definition = models.CollectionChan + " " + collectedDefinition(parseField(x.Elem()))

	// Function (without Receivers)
	// https://go.googlesource.com/example/+/HEAD/gotypes#function-and-method-types
	case *types.Signature:
		var definition strings.Builder

		// set the parameters.
		definition.WriteString(models.CollectionFunc + "(")
		for i := 0; i < x.Params().Len(); i++ {
			definition.WriteString(collectedDefinition(parseField(x.Params().At(i).Type())))
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
			definition.WriteString(collectedDefinition(parseField(x.Results().At(i).Type())))
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
				definition.WriteString(collectedDefinition(parseField(x.Method(i).Type())) + "; ")
			}

			for i := 0; i < x.NumEmbeddeds(); i++ {
				definition.WriteString(collectedDefinition(parseField(x.EmbeddedType(i))) + "; ")
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
			// a deepcopy of subfield is returned, then modified.
			subfield := parseField(x.Field(i).Type()).Deepcopy(nil)
			subfield.VariableName = "." + x.Field(i).Name()
			subfield.Name = x.Field(i).Name()
			setTags(subfield, x.Tag(i))
			subfield.Parent = field
			field.Fields = append(field.Fields, subfield)

			if x.Field(i).Embedded() {
				subfield.Embedded = true
			}

			definition.WriteString(subfield.Name + " " + subfield.FullDefinition() + "; ")

			// Due to the possibility of cyclic structs,
			// all subfields are deepcopied with len([]Fields) == (0:?).
			//
			// In order to correctly represent a deepcopied subfield,
			// point its fields back to the cached field []Fields,
			// which are eventually filled.
			//
			// cachedsubfield.Fields pointer is never modified.
			if cachedsubfield, ok := fieldcache[x.Field(i).String()]; ok {
				subfield.Fields = cachedsubfield.Fields
			}
		}
		definition.WriteString("}")
		field.Definition = definition.String()

	default:
		fmt.Printf("WARNING: could not parse type %v\n", x.String())
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

// collectedDefinition determines the full definition for a collected type in a collection.
//
// collectedDefinition can be called in the parser, but ONLY because collections are NOT cached.
func collectedDefinition(collected *models.Field) string {
	// a generated file's package == setup file's package.
	//
	// when the field is defined in the setup file (i.e `Collection`),
	// it is parsed with the setup file's package (i.e `copygen.Collection`).
	//
	// do NOT reference it by package in the generated file (i.e `Collection`).
	if collected.Import == setupPkgPath {
		return collected.Definition
	}

	// when a setup file imports the package it will output to,
	// do NOT reference the fields defined in the output package, by package.
	if outputPkgPath != "" && collected.Import == outputPkgPath {
		return collected.Definition
	}

	// when a field's import uses an alias, reassign the package reference.
	if aliasPkg, ok := aliasImportMap[collected.Import]; ok {
		return aliasPkg + "." + collected.Definition
	}

	return collected.FullDefinition()
}
