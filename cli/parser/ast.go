package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

// astLocateType finds the location of a file containing a type declaration in order to
// determine its import path, actual (non-aliased) package name, name, and definition.
func astLocateType(file *ast.File, imprt, name string) (string, string, string, string, error) {
	// traverse through the file's imports to determine the types actual package name.
	for _, importSpec := range file.Imports {
		importPath := importSpec.Path.Value
		if importPath == imprt {
			base := filepath.Base(importPath)
			// [:removes the last `"` from the package name]
			actualpkg := base[:len(base)-1]

			var definition string
			if actualpkg != "" {
				definition = actualpkg + "." + name
			} else {
				definition = name
			}
			return importPath, actualpkg, name, definition, nil
		}
	}
	return "", "", "", "", fmt.Errorf("Could not locate type %q with import %v.", name, imprt)
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
	return nil, fmt.Errorf("The type %q could not be found in the Abstract Syntax Tree.", typename)
}

// astSelectorSearch searches for a selector of a TypeSpec in an Abstract Syntax Tree.
func astSelectorSearch(ts *ast.TypeSpec, selector string) *ast.SelectorExpr {
	var astselector *ast.SelectorExpr
	ast.Inspect(ts, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.SelectorExpr:
			pkg := x.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
			name := x.Sel.Name           // 'Logger' in 'Field log.Logger'
			if pkg == "" && selector == name {
				astselector = x
			} else if selector == pkg+"."+name {
				astselector = x
			}
			return false
		default:
			return true
		}
	})
	return astselector
}

// parseASTFieldName parses an *ast.Field (node) for its package, name, definition, and pointer value.
func parseASTFieldName(field ast.Node) (string, string, string, string) {
	var pkg, name, def, ptr string
	ast.Inspect(field, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.SelectorExpr:
			// FieldInfo is always in a selector expression.
			pkg += x.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
			name += x.Sel.Name           // 'Logger' in 'Field log.Logger'
			return false
		case *ast.StarExpr:
			ptr += "*"
			return true
		default:
			return true
		}
	})
	if pkg != "" {
		def = pkg + "." + name
	} else {
		def = name
	}
	return pkg, name, def, ptr
}
