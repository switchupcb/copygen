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

	// check the cache for a parsed field.
	var cachekey string
	if x, ok := typ.(*types.Named); ok {

		// structs and interface `go/types` strings aren't unique enough,
		// so we must check for the type's import.
		var typeImport, typePackage string
		if x.Obj().Pkg() != nil {
			typeImport = x.Obj().Pkg().Path()

			if _, ok := aliasImportMap[typeImport]; ok {
				typePackage = aliasImportMap[typeImport]
			} else {
				typePackage = x.Obj().Pkg().Name()
			}
		}

		cachekey = typeImport + typePackage + typ.String()
	} else {
		cachekey = typ.String()
	}

	if cached, ok := fieldcache[cachekey]; ok {
		return cached
	}

	// build the field in the cache.
	cachefield := new(models.Field)

	// do NOT cache pointers.
	if typ.String()[0:1] != models.Pointer {
		fieldcache[cachekey] = cachefield
	}

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
		cachefield.Definition = x.Obj().Name()
		setFieldImportAndPackage(cachefield, x.Obj().Pkg())

		// Struct Types
		// https://go.googlesource.com/example/+/HEAD/gotypes#struct-types
		if s, ok := x.Underlying().(*types.Struct); ok {
			parseStructField(cachefield, s)
		}

	// Basic Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#basic-types
	case *types.Basic:
		cachefield.Definition = x.Name()

	// Simple Composite Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#simple-composite-types
	case *types.Pointer:
		// underlyingfield is the cache representation of the underlying field (i.e `int`, `*int`).
		underlyingfield := parseField(x.Elem())
		cachefield.Import = underlyingfield.Import
		cachefield.Package = underlyingfield.Package
		cachefield.Fields = underlyingfield.Fields

		// set the definition accordingly (i.e `*int`, `**int`).
		cachefield.Pointer = models.Pointer
		if underlyingfield.IsPointer() {
			cachefield.Definition = models.CollectionPointer + underlyingfield.Definition
		} else {
			cachefield.Definition = underlyingfield.Definition
		}

	case *types.Array:
		underlyingfield := parseField(x.Elem())
		cachefield.Definition = "[" + fmt.Sprint(x.Len()) + "]" + underlyingfield.Definition

	case *types.Slice:
		underlyingfield := parseField(x.Elem())
		cachefield.Definition = models.CollectionSlice + underlyingfield.Definition

	case *types.Map:
		keyfield := parseField(x.Key())
		valfield := parseField(x.Elem())
		cachefield.Definition = models.CollectionMap + "[" + keyfield.Definition + "]" + valfield.Definition

	case *types.Chan:
		underlyingfield := parseField(x.Elem())
		cachefield.Definition = models.CollectionChan + " " + underlyingfield.Definition

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
		hasResults := x.Results().Len() > 1
		if hasResults {
			definition.WriteString("(")
		}
		for i := 0; i < x.Results().Len(); i++ {
			definition.WriteString(parseField(x.Results().At(i).Type()).Definition)
			if i+1 != x.Results().Len() {
				definition.WriteString(", ")
			}
		}
		if hasResults {
			definition.WriteString(")")
		}

		cachefield.Definition = definition.String()

	// Interface Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#interface-types
	case *types.Interface:
		if x.Empty() {
			cachefield.Definition = x.String()
		} else {
			var definition strings.Builder
			definition.WriteString(models.CollectionInterface + "{")
			for i := 0; i < x.NumMethods(); i++ {
				definition.WriteString(parseField(x.Method(i).Type()).Definition + "; ")
			}
			definition.WriteString("}")
			cachefield.Definition = definition.String()
		}

	// Unnamed Struct (struct{})
	case *types.Struct:
		cachefield.Definition = "struct{}"

	default:
		fmt.Printf("WARNING: could not parse type %v\n", x.String())
	}

	cachefield.VariableName = "." + alphastring(cachefield.Definition)
	return cachefield
}

// parseStructField parses a struct field.
func parseStructField(field *models.Field, x *types.Struct) {
	for i := 0; i < x.NumFields(); i++ {
		// a deepcopy of subfield is returned and modified.
		subfield := parseField(x.Field(i).Type()).Deepcopy(nil)
		subfield.VariableName = "." + x.Field(i).Name()
		subfield.Name = x.Field(i).Name()
		setTags(subfield, x.Tag(i))
		subfield.Parent = field
		field.Fields = append(field.Fields, subfield)

		// all subfields are deepcopied with Fields[0].
		//
		// in order to correctly represent a deepcopied struct field,
		// we must point its fields back to the cached field.Fields,
		// which will eventually be filled.
		//
		// cachedsubfield.Fields are never modified.
		if cachedsubfield, ok := fieldcache[subfield.Import+subfield.Package+x.Field(i).String()]; ok {
			subfield.Fields = cachedsubfield.Fields
		}
	}
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
		fmt.Printf("WARNING: could not parse tag for field %v\n%v", field.FullName(""), err)
	}

	field.Tags = make(map[string]map[string][]string, tags.Len())
	for _, tag := range tags.Tags() {
		field.Tags[tag.Key] = map[string][]string{
			tag.Name: tag.Options,
		}
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
