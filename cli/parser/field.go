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

// SearchForField searches for an *ast.Field which is parsed into a field model.
//
// The field search process involves a FieldSearcher that sets up and executes a field search in order to load field data.
// In the context of the program, a FieldSearcher's fieldSearch contains a TypeField as its top-level Searcher.
// The original setup file (i.e setup.go) is used to locate a field's actual import and package.
// Then, the files that compose this package are searched for the declaration of the field containing its data and sub-fields.
func (fs *FieldSearcher) SearchForField(setupfile *ast.File, setimport, setpkg, setname string) fieldSearch {
	if fs.cache == nil {
		fs.cache = make(map[string]fieldSearch)
	}
	if cachedsearch, ok := fs.cache[setimport+setpkg+setname]; ok {
		return cachedsearch
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
	actualimport, actualpkg, actualname, definition, err := astLocateType(setupfile, setimport, setname)
	if err != nil {
		fs.FieldSearch = fieldSearch{Error: err}
		return fs.FieldSearch
	}

	// set up the field searcher (and set the properties of the field) using data from the actual import file.
	fs.FieldSearch.Error = fs.FieldSearch.setup(actualimport, actualpkg, actualname, definition)
	if fs.FieldSearch.Error != nil {
		return fs.FieldSearch
	}
	// A top-level searcher returns itself.
	if fs.FieldSearch.Searcher.Parent == nil {
		fs.FieldSearch.Result = fs.FieldSearch.Searcher
	}

	// exeute the search
	fs.FieldSearch.Result, fs.FieldSearch.Error = fs.execute()
	if fs.FieldSearch.Error != nil {
		fs.FieldSearch.Error = fmt.Errorf("An error occurred while searching for the Field %q with import: %v.\n%v", fs.FieldSearch.Searcher.FullName(""), fs.FieldSearch.Searcher.Import, fs.FieldSearch.Error)
		return fs.FieldSearch
	}
	return fs.FieldSearch
}

// FieldSearcher represents a searcher that uses Abstract Syntax Tree analysis to find fields of a type.
type FieldSearcher struct {
	// The current search for the field searcher; or nil.
	FieldSearch fieldSearch

	// The options applied to fields during a search.
	// map[option][]values
	Options map[string][]string

	// A key value cache used to reduce the amount of AST operations.
	cache map[string]fieldSearch
}

// fieldSearch represents a search for a field.
type fieldSearch struct {
	// The searcher that initiated the field search; or nil.
	Searcher *models.Field

	// The typespec of the searcher that initiated the field search.
	SearcherTypeSpec *ast.TypeSpec

	// The files discovered during the search.
	Files []*ast.File

	// The file that holds the type declaration for the searcher.
	DecFile *ast.File

	// The types info for the search.
	Info types.Info

	// The resulting field found by the search.
	// Result can only contain custom fields OR a basic field; NOT both.
	Result *models.Field

	// The error that occurred during the search; or nil.
	Error error

	// Whether the results contain a basic field.
	// There can only ever be one basic field in a search (since a basic type doesn't contain other fields).
	isBasic bool

	// The current depth-level of the fieldSearch.
	Depth int

	// The maximum allowed depth-level of the fieldSearch.
	MaxDepth int
}

// setup sets up a field search for execution by checking the types of an *ast.Fileset (with *ast.Files)
// and loading types.Info and an *ast.TypeSpec into the search.
func (fs *fieldSearch) setup(imprt, pkg, name, def string) error {
	if fs.Searcher == nil {
		fs.Searcher = &models.Field{
			Import:     imprt,
			Package:    pkg,
			Name:       name,
			Definition: def,
		}
	} else {
		fs.Searcher.Import = imprt
		fs.Searcher.Package = pkg
		fs.Searcher.VariableName = "." + name
		fs.Searcher.Name = name
		fs.Searcher.Definition = def
	}

	// TODO: insert field options

	packages, err := packages.Load(&packages.Config{Logf: nil}, fs.Searcher.Import[1:len(fs.Searcher.Import)-1])
	if err != nil {
		return fmt.Errorf("An error occurred retrieving a package from the GOPATH: %v\n%v", fs.Searcher.Import, err)
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
		fs.Files = append(fs.Files, file)
	}

	// determine the types present in the package
	conf := types.Config{Importer: importer.Default()}
	fs.Info = types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err = (conf.Check(fs.Searcher.Package, fileset, fs.Files, &fs.Info))
	if err != nil {
		return fmt.Errorf("An error occurred determining the types of a package.\n%v", err)
	}

	// determine the TypeSpec for this search
	// find the actual file location of the field's type declaration using the setup file.
	var ts *ast.TypeSpec
	for _, file := range fs.Files {
		ts, _ = astTypeSearch(file, fs.Searcher.Name)
		if ts != nil {
			fs.DecFile = file
			break
		}
	}
	if ts == nil {
		return fmt.Errorf("The type %v could not be found in the AST for Field in package %q with import %v.\nIs the package up to date?", name, pkg, imprt)
	}
	fs.SearcherTypeSpec = ts
	return nil
}

// execute searches for a field's data (import, package, name, definition...) in a list of files' Abstract Syntax Tree.
func (fs *FieldSearcher) execute() (*models.Field, error) {
	fieldsearch := fs.FieldSearch.astFieldSearch(fs.Options, fs.cache)
	if fieldsearch.Error != nil {
		return nil, fieldsearch.Error
	}

	if fieldsearch.Result != nil || fieldsearch.isBasic {
		fs.cache[fs.FieldSearch.Searcher.Import+fs.FieldSearch.Searcher.Package+fs.FieldSearch.Searcher.Name] = fieldsearch
		return fieldsearch.Result, nil
	}
	return nil, fmt.Errorf("The type %v could not be loaded from the specified module: %v", fs.FieldSearch.Searcher, fs.FieldSearch.Searcher.Import)
}

// astFieldSearch searches through an ast.Typespec for fields.
func (fs fieldSearch) astFieldSearch(options map[string][]string, cache map[string]fieldSearch) fieldSearch {
	// find the fields of the main field if the max depth-level has not been reached.
	/* if fs.Depth < fs.MaxDepth { switch...}  */
	fs.Depth++
	switch x := fs.Info.Types[fs.SearcherTypeSpec.Type].Type.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xField := x.Field(i)
			subfield := &models.Field{
				Parent:       fs.Searcher,
				VariableName: "." + xField.Name(),
				Name:         xField.Name(),
				Definition:   xField.Type().String(),
			}
			if !isBasic(xField.Type()) {
				newFieldSearcher := FieldSearcher{Options: options, cache: cache}
				newFieldSearcher.FieldSearch.Searcher = subfield

				// find the custom type field.
				var err error
				subfield.Import, subfield.Package, subfield.Name, subfield.Definition, subfield.Pointer, err = astSubFieldSearch(fs.DecFile, subfield.Parent.Import, subfield.Parent.Package, subfield.Name, subfield.Definition)
				if err != nil {
					return fieldSearch{Error: err}
				}
				newfieldsearch := newFieldSearcher.SearchForField(fs.DecFile, subfield.Import, subfield.Package, subfield.Name)
				if newfieldsearch.Error != nil {
					return newfieldsearch // doesn't return the top-level search so it's expected that an error is always checked for.
				}
				fs.Result.Fields = append(fs.Result.Fields, newfieldsearch.Result)
			} else {
				fs.Searcher.Fields = append(fs.Searcher.Fields, subfield)
			}
		}
		return fs
	// interfaces have method fields
	case *types.Interface:
		for i := 0; i < x.NumMethods(); i++ {
			xMethod := x.Method(i)
			// interface functions are declared in the same scope as the interface
			subfield := &models.Field{
				VariableName: "." + xMethod.Name() + "(%)",
				Import:       fs.Searcher.Import,
				Package:      fs.Searcher.Package,
				Name:         xMethod.Name(),
				Definition:   xMethod.Type().String(),
				Parent:       fs.Searcher,
			}
			fs.Result.Fields = append(fs.Result.Fields, subfield)
		}
		fs.isBasic = true
		return fs
	// if no fields are present, this is a basic type.
	default:
		fs.isBasic = true
		return fs
	}
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

// // parseFieldOptions sets the options of a field.
// func parseFieldOptions(options map[string][]string) {
// 	return ""
// }
