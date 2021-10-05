package parser

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
	"golang.org/x/tools/go/packages"
)

// SearchForTypeField searches for an *ast.Field which is parsed into a (type) field model.
//
// The field search process involves a FieldSearcher that sets up and executes a field search in order to load field data.
// In the context of the program, a top-level field with no parents is a TypeField.
// The original setup file (i.e setup.go) is used to locate a field's actual import and package.
// Then, the files that compose this package are searched for the declaration of the field containing its data and sub-fields.
func (fs *FieldSearcher) SearchForTypeField(setupfile *ast.File, setimport, setpkg, setname string) (*models.Field, error) {
	if fs.cache == nil {
		fs.cache = make(map[string]*models.Field)
	}

	// There is a difference between the parameterized properties (which are parsed from a "setup file")
	// and the actual file properties (which are parsed from the file containing the field's type declaration).
	//
	// SearchForField is passed a file, modularized import path, aliased package, and (type) name from a setup file.
	// This means that imports with aliased packages (i.e c "github.com/switchupb/copygen/examples/main/converter")
	// will be parsed from the Copygen interface function. However, a module import != importpath && alias != package.
	// In order to solve this, we must locate the types ACTUAL properties from its declaration in antoher file.
	//
	// find the actual file location of the field's type declaration using the setup file.
	actualimport, err := astLocateImport(setupfile, setimport, setpkg, setname)
	if err != nil {
		return nil, err
	}

	// set up and execute the field searcher using data from the actual import file.
	def := setpkg
	if def != "" {
		def += "."
	}
	def += setname

	// when fs.Field == nil; a TypeField is instantiated
	fs.Field = nil
	fs.SearchInfo = FieldSearchInfo{Depth: 0, MaxDepth: 0}
	if err := fs.execute(actualimport, setpkg, setname, def); err != nil {
		return nil, fmt.Errorf("An error occurred while searching for the Field %q of package %q with import: %v.\n%v", setname, setpkg, setimport, err)

	}
	return fs.Field, nil
}

// FieldSearcher represents a searcher that uses Abstract Syntax Tree analysis to find fields of a type.
type FieldSearcher struct {
	// The field that initiates the search.
	Field *models.Field

	// The current search information for the field searcher.
	SearchInfo FieldSearchInfo

	// The options that pertain to a field (and its subfields).
	Options []Option

	// A key value cache used to reduce the amount of AST operations.
	cache map[string]*models.Field
}

// FieldSearchInfo represents the info for a field search.
type FieldSearchInfo struct {
	// The typespec of the searcher that initiated the field search.
	SearcherTypeSpec *ast.TypeSpec

	// The files discovered during the search.
	Files []*ast.File

	// The file that holds the type declaration for the searcher.
	DecFile *ast.File

	// The types info for the search.
	Info types.Info

	// Whether the results contain a basic field.
	// There can only ever be one basic field in a search (since a basic type doesn't contain other fields).
	isBasic bool

	// The current depth-level of the fieldSearch.
	Depth int

	// The maximum allowed depth-level of the fieldSearch.
	MaxDepth int
}

// execute runs a field search by checking the types of an *ast.Fileset (with *ast.Files), loading types.Info and an *ast.TypeSpec
// then searching for a field and it's subfields.
func (fs *FieldSearcher) execute(imprt, pkg, name, def string) error {
	if cachedsearch, ok := fs.cache[imprt+pkg+name]; ok {
		fs.Field = cachedsearch
		return nil
	}

	// setup the field
	// if the field is nil, it's a TypeField
	if fs.Field == nil {
		fs.Field = &models.Field{
			Import:     imprt,
			Package:    pkg,
			Name:       name,
			Definition: def,
		}
	} else {
		fs.Field.Import = imprt
		fs.Field.Package = pkg
		fs.Field.Name = name
		fs.Field.Definition = def
		fs.Field.VariableName = "." + name
	}
	setFieldOptions(fs.Field, fs.Options)
	fs.SearchInfo.MaxDepth += fs.Field.Options.Depth

	// load the package the field is located in
	packages, err := packages.Load(&packages.Config{Logf: nil}, imprt[1:len(imprt)-1])
	if err != nil {
		return fmt.Errorf("An error occurred retrieving a package from the GOPATH: %v\n%v", imprt, err)
	}
	var gofiles []string
	for _, pkg := range packages {
		gofiles = append(gofiles, pkg.GoFiles...)
	}

	// prepare the loaded package for AST analysis
	fileset := token.NewFileSet()
	for _, filepath := range gofiles {
		file, err := parser.ParseFile(fileset, filepath, nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("An error occurred parsing a file for the matcher: %v\n%v", filepath, err)
		}
		fs.SearchInfo.Files = append(fs.SearchInfo.Files, file)
	}

	// determine the types present in the package
	conf := types.Config{Importer: importer.Default()}
	fs.SearchInfo.Info = types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err = (conf.Check(pkg, fileset, fs.SearchInfo.Files, &fs.SearchInfo.Info))
	if err != nil {
		return fmt.Errorf("An error occurred determining the types of a package.\n%v", err)
	}

	// determine the TypeSpec for this search using the actual typename (i.e `DomainUser` in `User DomainUser`)
	var typename string
	splitdef := strings.Split(def, ".")
	if len(splitdef) == 1 {
		typename = def
	} else {
		typename = splitdef[1]
	}
	var ts *ast.TypeSpec
	for _, file := range fs.SearchInfo.Files {
		ts, _ = astTypeSearch(file, typename)
		if ts != nil {
			fs.SearchInfo.DecFile = file
			break
		}
	}
	if ts == nil {
		return fmt.Errorf("The type declaration for the Field %q with import %v could not be found in the AST.\nIs the package up to date?", fs.Field.FullName(""), imprt)
	}
	fs.SearchInfo.SearcherTypeSpec = ts

	// find the fields of the main field if the max depth-level has not been reached.
	var subfields []*models.Field
	subfields, err = fs.astFieldSearch()
	if err != nil {
		return err
	}

	fs.Field.Fields = subfields
	fs.cache[fs.Field.Import+fs.Field.Package+fs.Field.Name] = fs.Field
	return nil
}

// astFieldSearch searches through an ast.Typespec for sub-fields.
func (fs *FieldSearcher) astFieldSearch() ([]*models.Field, error) {
	var subfields []*models.Field
	switch x := fs.SearchInfo.Info.Types[fs.SearchInfo.SearcherTypeSpec.Type].Type.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xField := x.Field(i)

			// create a new typefield if a subfield is a custom type
			if (fs.SearchInfo.MaxDepth == 0 || fs.SearchInfo.Depth < fs.SearchInfo.MaxDepth) && !isBasic(xField.Type()) {
				// find the actual custom type field info
				splitdefinition := strings.Split(xField.Type().String(), ".")
				defPkg := splitdefinition[0] // i.e `log` in `log.Logger`
				if defPkg == "" {
					defPkg = fs.Field.Package
				}
				actualimport, err := astLocateImport(fs.SearchInfo.DecFile, fs.Field.Import, defPkg, xField.Name())
				if err != nil {
					return nil, fmt.Errorf("An error occurred searching for subfield %q of type %q\n%v", fs.Field.FullName(""), xField.Type().String(), err)
				}

				// a newFieldSearcher contains the same options and cache, but new field search info
				newFieldSearcher := FieldSearcher{
					SearchInfo: FieldSearchInfo{
						Depth:    fs.SearchInfo.Depth + 1,
						MaxDepth: fs.SearchInfo.MaxDepth,
					},
					Options: fs.Options,
					cache:   fs.cache,
				}

				// Ensure a new TypeField is NOT created.
				newFieldSearcher.Field = &models.Field{Parent: fs.Field, Definition: xField.Type().String()}

				// Search for the subfields of the subfield
				if err := newFieldSearcher.execute(actualimport, defPkg, xField.Name(), xField.Type().String()); err != nil {
					return nil, err
				}
				subfields = append(subfields, newFieldSearcher.Field)
			} else {
				subfield := &models.Field{
					Parent:       fs.Field,
					VariableName: "." + xField.Name(),
					Name:         xField.Name(),
					Definition:   xField.Type().String(),
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
				Import:       fs.Field.Import,
				Package:      fs.Field.Package,
				Name:         xMethod.Name(),
				Definition:   xMethod.Type().String(),
				Parent:       fs.Field,
			}
			setFieldOptions(subfield, fs.Options)
			subfields = append(subfields, subfield)
		}
	default:
		// if no fields are present, this is a basic type.
	}
	return subfields, nil
}

// astSubFieldSearch searches through an AST using limited information to determine
// an import, package, name, and definition and pointer (if applicable) in a setup file.
func astSubFieldSearch(file *ast.File, parentImport, parentPkg, typeName, definition string) (string, string, string, string, string, error) {
	var imprt, pkg, name, def, ptr string
	splitDefinition := strings.Split(definition, ".")
	if len(splitDefinition) >= 2 {
		definitionPkg := splitDefinition[0]  // 'log' in 'Field log.Logger'
		definitionName := splitDefinition[1] // 'Logger' in 'Field log.Logger'

		// use the selector on a custom type from a different package to determine its field.
		if definitionPkg != parentPkg {
			// find the type in the AST
			ts, err := astTypeSearch(file, name)
			if err == nil {
				return "", "", "", "", "", fmt.Errorf("The type %v.%v could not be found in the AST. Is the package up to date?", pkg, name)
			}
			sel := astSelectorSearch(ts, definitionPkg+"."+definitionName)
			pkg, name, def, ptr = parseASTFieldName(sel)
		} else {
			imprt = parentImport
			pkg = parentPkg
		}
		def = definitionPkg + "." + definitionName
	}
	return imprt, pkg, name, def, ptr, nil
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
		if options[i].Category == "convert" && options[i].Regex[1].MatchString(field.FullName("")) {
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
		if options[i].Category == "deepcopy" && options[i].Regex[0].MatchString(field.FullName("")) {
			field.Options.Deepcopy = true
			break
		}
	}
}

// setDepthOption sets a field's depth option.
func setDepthOption(field *models.Field, options []Option) {
	// A depth option can only be set to a field once, so use the last one
	for i := len(options) - 1; i > -1; i-- {
		if options[i].Category == "depth" && options[i].Regex[0].MatchString(field.FullName("")) {
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
		if options[i].Category == "map" && options[i].Regex[0].MatchString(field.FullName("")) {
			if value, ok := options[i].Value.(string); ok {
				field.Options.Map = value
				break
			}
		}
	}
}
