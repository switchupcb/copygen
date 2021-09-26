package loader

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

// AST provides AST analysis to find fields.
type AST struct {
	cache map[string][]models.Field // ASTCache is a key value cache used to reduce the amount of AST operations.
}

// ASTSearch searches a .go source file for a type and its fields.
func (a *AST) ASTSearch(imprt string, pkg string, typename string) ([]models.Field, error) {
	if a.cache == nil {
		a.cache = make(map[string][]models.Field)
	}
	if search, ok := a.cache[imprt+pkg+typename]; ok {
		return search, nil
	}

	packages, err := packages.Load(&packages.Config{Logf: nil}, imprt)
	if err != nil {
		return nil, fmt.Errorf("An error occurred retrieving a package from the GOPATH: %v\n%v", imprt, err)
	}
	var gofiles []string
	for _, pkgs := range packages {
		gofiles = append(gofiles, pkgs.GoFiles...)
	}

	fileset := token.NewFileSet()
	var astFiles []*ast.File
	for _, filepath := range gofiles {
		file, err := parser.ParseFile(fileset, filepath, nil, parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("An error occurred parsing a file for the matcher: %v\n%v", filepath, err)
		}
		astFiles = append(astFiles, file)
	}

	// check the package types
	conf := types.Config{Importer: importer.Default()}
	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	_, err = (conf.Check(pkg, fileset, astFiles, &info))
	if err != nil {
		return nil, fmt.Errorf("An error occurred determining the types of a package.\n%v", err)
	}

	// find the type in the AST
	var ts *ast.TypeSpec
	for _, file := range astFiles {
		ts, _ = astTypeSearch(file, typename)
		if ts != nil {
			break
		}
	}
	if ts == nil {
		return nil, fmt.Errorf("The type %v.%v could not be found in the AST. Is the package up to date?", pkg, typename)
	}

	// find the fields
	for _, file := range astFiles {
		fieldSearch := a.astFieldSearch(info, file, ts, imprt, pkg)
		if fieldSearch.Error != nil {
			return nil, fieldSearch.Error
		}

		if len(fieldSearch.Fields) != 0 || fieldSearch.Basic {
			a.cache[imprt+pkg+typename] = fieldSearch.Fields
			return fieldSearch.Fields, nil
		}
	}
	return nil, fmt.Errorf("The type %v could not be loaded from the specified module: %v\n", pkg+"."+typename, imprt)
}

// fieldSearch represents a search for a field.
type fieldSearch struct {
	Fields []models.Field // The fields present in the search.
	Basic  bool           // Whether there are fields that are basic.
	Error  error          // Whether an error occured.
}

// astFieldSearch searches through an ast.Typespec for fields.
func (a *AST) astFieldSearch(info types.Info, file *ast.File, ts *ast.TypeSpec, imprt string, pkg string) fieldSearch {
	var fields []models.Field
	switch x := info.Types[ts.Type].Type.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xField := x.Field(i)
			fieldname := xField.Name()
			definition := xField.Type().String()
			field := models.Field{
				Name:       fieldname,
				Definition: definition,
			}

			// if a field is a custom type it may have fields of its own
			if !isBasic(xField.Type()) {
				// find the custom type field.
				splitDefinition := strings.Split(field.Definition, ".")
				if len(splitDefinition) == 2 {
					definitionPkg := splitDefinition[0]
					definitionType := splitDefinition[1]

					// use the selector on a custom type to determine its field
					var newImprt, newPkg, newType string
					if definitionPkg != pkg {
						sel := astSelectorSearch(ts, definitionPkg+"."+definitionType)
						if sel == nil {
							return fieldSearch{
								Error: fmt.Errorf("Could not find the selector for the %v in-depth field %v", ts.Name.Name, field.Definition),
							}
						}
						newImprt, newPkg, newType = astLocateType(file, sel)
					} else {
						newImprt = imprt
						newPkg = pkg
						newType = definitionType
					}
					depthFields, err := a.ASTSearch(newImprt, newPkg, newType)
					if err != nil {
						fmt.Printf("WARNING: An error occurred searching for the %v in-depth field '%v' with import \"%v\".\n%v\n", newType, field.Definition, newImprt, err)
					}
					field.Fields = depthFields
				}
			}
			fields = append(fields, field)
		}
	// interfaces have method fields
	case *types.Interface:
		for i := 0; i < x.NumMethods(); i++ {
			xMethod := x.Method(i)
			fieldname := xMethod.Name()
			definition := xMethod.Type().String()
			field := models.Field{
				Name:       fieldname,
				Definition: definition,
			}
			fields = append(fields, field)
		}
	// if no fields are present, this is a basic type.
	default:
		return fieldSearch{
			Basic: true,
		}
	}
	return fieldSearch{
		Fields: fields,
	}
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

// astTypeSearch searches through an ast.File for ast.Types.
func astTypeSearch(file *ast.File, typename string) (*ast.TypeSpec, error) {
	for _, decl := range file.Decls {
		if gendecl, ok := decl.(*ast.GenDecl); ok {
			if gendecl.Tok == token.TYPE {
				for _, spec := range gendecl.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if typename == ts.Name.Name {
							return ts, nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("The type %v could not be found in the AST.", typename)
}

// astDepthFields finds fields of fields using an AST.
func astSelectorSearch(ts *ast.TypeSpec, selector string) *ast.SelectorExpr {
	if strct, ok := ts.Type.(*ast.StructType); ok {
		for _, field := range strct.Fields.List {
			if sel, ok := field.Type.(*ast.SelectorExpr); ok {
				fieldTypePkg := sel.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
				fieldTypeName := sel.Sel.Name           // 'Logger' in 'Field log.Logger'
				if selector == fieldTypePkg+"."+fieldTypeName {
					return sel
				}
			}
		}
	}
	return nil
}

// astLocateType finds the import path, package, and typename of a type in an AST.
func astLocateType(file *ast.File, sel *ast.SelectorExpr) (string, string, string) {
	fieldTypePkg := sel.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
	fieldTypeName := sel.Sel.Name           // 'Logger' in 'Field log.Logger'

	// don't alter the original file's slice
	var checkedImprts []*ast.ImportSpec
	for _, imprt := range file.Imports {
		checkedImprts = append(checkedImprts, imprt)
	}

	// check imports that have variable names
	for i := len(checkedImprts) - 1; i >= 0; i-- {
		imprt := checkedImprts[i]

		// if an import has a variable name
		if imprt.Name != nil {
			// if an import variable matches the package name (i.e 'log' in 'log.Logger')
			if fieldTypePkg == imprt.Name.Name {
				return imprt.Path.Value, fieldTypePkg, fieldTypeName
			} else {
				// remove
				checkedImprts = checkedImprts[:len(checkedImprts)-1]
			}
		}
	}

	// check remaining imports (that don't have variable names)
	for _, imprt := range checkedImprts {
		imprtPath := imprt.Path.Value
		imprtPath = imprtPath[1 : len(imprtPath)-1] // "log" -> log
		if fieldTypePkg == imprtPath {
			return imprtPath, fieldTypePkg, fieldTypeName
		}
	}
	return "", fieldTypePkg, fieldTypeName
}

// printFields shows a tree of fields for a given type.
func PrintFields(typename string, fields []models.Field, tabs string) {
	if tabs == "" {
		fmt.Println(tabs + "type " + typename)
	}

	tabs += "\t" // field tab
	for _, field := range fields {
		fmt.Println(tabs + field.Name + "\t" + field.Definition)
		if len(field.Fields) != 0 {
			PrintFields(field.Definition, field.Fields, tabs+"\t")
		}
	}
}
