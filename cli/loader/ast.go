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

// ASTSearch searches a .go source file for a type and its fields.
func ASTSearch(imprt string, pkg string, typename string) ([]models.Field, error) {
	packages, err := packages.Load(&packages.Config{Logf: nil}, imprt)
	if err != nil {
		return nil, fmt.Errorf("There was an error retrieving a package from the GOPATH: %v\n%v", imprt, err)
	}
	var gofiles []string
	for _, pkgs := range packages {
		gofiles = append(gofiles, pkgs.GoFiles...)
	}

	fs := token.NewFileSet()
	var astfs []*ast.File
	for _, filepath := range gofiles {
		file, err := parser.ParseFile(fs, filepath, nil, parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("An error occured parsing the file: %v\n%v", filepath, err)
		}
		astfs = append(astfs, file)
	}

	// check the package types
	conf := types.Config{Importer: importer.Default()}
	info := types.Info{Types: make(map[ast.Expr]types.TypeAndValue), Defs: make(map[*ast.Ident]types.Object)}
	_, err = (conf.Check(pkg, fs, astfs, &info))
	if err != nil {
		return nil, fmt.Errorf("An error occured determining the types of a package.\n%v", err)
	}

	// find the type in the AST
	var ts *ast.TypeSpec
	for _, f := range astfs {
		ts, _ = astTypeSearch(f, typename)
		if ts != nil {
			break
		}
	}
	if ts == nil {
		return nil, fmt.Errorf("The type %v.%v could not be found in the AST. Is the package up to date?", pkg, typename)
	}

	// find the fields
	for _, f := range astfs {
		fields := astFieldSearch(info, f, ts, imprt, pkg)
		if fields.Error != nil {
			return nil, fields.Error
		}

		if len(fields.Fields) != 0 || fields.Basic {
			return fields.Fields, nil
		}
	}
	return nil, fmt.Errorf("The type %v could not be loaded from the specified module: %v\n", pkg+"."+typename, imprt)
}

// astTypeSearch searches through an ast.File for ast.Types.
func astTypeSearch(f *ast.File, t string) (*ast.TypeSpec, error) {
	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.TYPE {
				for _, s := range gd.Specs {
					if ts, ok := s.(*ast.TypeSpec); ok {
						if t == ts.Name.Name {
							return ts, nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("The type %v could not be found in the AST.", t)
}

// astDepthFields finds fields of fields using an AST.
func astSelectorSearch(f *ast.File, ts *ast.TypeSpec, selector string) *ast.SelectorExpr {
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
func astLocateType(f *ast.File, fld *ast.SelectorExpr) (string, string, string) {
	fldTypePkg := fld.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
	fldTypeName := fld.Sel.Name           // 'Logger' in 'Field log.Logger'

	// don't alter the original file's slice
	var checkedImprts []*ast.ImportSpec
	for _, v := range f.Imports {
		checkedImprts = append(checkedImprts, v)
	}

	// check imports that have variable names
	for i := len(checkedImprts) - 1; i >= 0; i-- {
		imprt := checkedImprts[i]

		// if an import has a variable name
		if imprt.Name != nil {
			// if an import variable matches the package name (i.e 'log' in 'log.Logger')
			if fldTypePkg == imprt.Name.Name {
				return imprt.Path.Value, fldTypePkg, fldTypeName
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
		if fldTypePkg == imprtPath {
			return imprtPath, fldTypePkg, fldTypeName
		}
	}
	return "", fldTypePkg, fldTypeName
}

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

// FieldSearch represents a search for a field.
type fieldSearch struct {
	Fields []models.Field // The fields present in the search.
	Basic  bool           // Whether there are fields are basic.
	Error  error          // Whether an error occured.
}

// astFieldSearch searches through an ast.Typespec for fields.
func astFieldSearch(info types.Info, f *ast.File, ts *ast.TypeSpec, imprt string, pkg string) fieldSearch {
	var fields []models.Field
	switch x := info.Types[ts.Type].Type.(type) {
	// structs have fields that can have fields.
	case *types.Struct:
		for i := 0; i < x.NumFields(); i++ {
			xfield := x.Field(i)
			fieldname := xfield.Name()
			definition := xfield.Type().String()
			field := models.Field{
				Name:       fieldname,
				Definition: definition,
			}

			// if a field is a custom type it may have fields of its own
			if !isBasic(xfield.Type()) {
				// find the custom type field.
				defs := strings.Split(field.Definition, ".")
				if len(defs) == 2 {
					dpkg := defs[0]
					dtyp := defs[1]

					// use the selector on a custom type to determine its field
					var nimprt, npkg, ntype string
					if dpkg != pkg {
						sel := astSelectorSearch(f, ts, dpkg+"."+dtyp)
						if sel == nil {
							return fieldSearch{
								Error: fmt.Errorf("Could not find the selector for the %v in-depth field %v", ts.Name.Name, field.Definition),
							}
						}
						nimprt, npkg, ntype = astLocateType(f, sel)
					} else {
						nimprt = imprt
						npkg = pkg
						ntype = dtyp
					}
					depthFields, err := ASTSearch(nimprt, npkg, ntype)
					if err != nil {
						fmt.Printf("WARNING: An error occured searching for the %v in-depth field '%v' with import \"%v\".\n%v\n", ntype, field.Definition, imprt, err)
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

// printFields shows a tree of fields for a given type.
func PrintFields(t string, fields []models.Field, tabs string) {
	if tabs == "" {
		fmt.Println(tabs + "type " + t)
	}

	tabs += "\t" // field tab
	for _, field := range fields {
		fmt.Println(tabs + field.Name + "\t" + field.Definition)
		if len(field.Fields) != 0 {
			PrintFields(field.Definition, field.Fields, tabs+"\t")
		}
	}
}
