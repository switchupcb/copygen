package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// splitASTComment splits an *ast.Comment for options parsing without the //.
func splitASTComment(c *ast.Comment, separator string) []string {
	// Remove the `//` and normalize space.
	comment := strings.TrimSpace(c.Text[2:])

	// Separate comment into 1 (option) + n parts
	return strings.Split(comment, separator)
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

// astLocateImport finds the actual import of a given package in a .go file.
// The import is used to load packages prior to a field search.
func astLocateImport(file *ast.File, fileImport, pkg, name string) (string, error) {
	// A type with no referenced package is declared in the same file.
	if pkg == "" {
		return fileImport, nil
	}

	// check the current file
	base := filepath.Base(fileImport)
	if pkg == base[:len(base)-1] {
		return fileImport, nil
	}

	for _, importSpec := range file.Imports {
		importPath := importSpec.Path.Value

		// check aliased imports (i.e `c "strconv"`)
		if importSpec.Name != nil && pkg == importSpec.Name.Name {
			return importPath, nil
		}

		// check stdlib imports (i.e `"log"`, `"strconv"`)
		if pkg == importPath[1:len(importPath)-1] {
			return importPath, nil
		}

		// check file imports (i.e `"github.com/switchupcb/copygen/models`)
		base := filepath.Base(importPath)
		if pkg == base[:len(base)-1] {
			return importPath, nil
		}
	}
	return "", fmt.Errorf("Could not locate type %q in file import %v.", pkg+" "+name, fileImport)
}

// astLocateType uses the location of a file containing a type declaration in order to
// determine the type's actual import path, package, typename, and definition.
func astLocateType(file *ast.File, sel *ast.SelectorExpr) (string, string, string, string, error) {
	typefieldPkg := sel.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
	typefieldName := sel.Sel.Name           // 'Logger' in 'Field log.Logger'
	definition := typefieldPkg
	if definition != "" {
		definition += "."
	}
	definition = typefieldName

	// traverse through the file's imports to determine the types actual package name.
	for _, importSpec := range file.Imports {
		importPath := importSpec.Path.Value

		// check imports that have variables names
		if importSpec.Name != nil {
			// if an import variable matches the package name (i.e 'log' in 'log.Logger')
			if typefieldPkg == importSpec.Name.Name {
				return importPath, typefieldPkg, typefieldName, definition, nil
			}
		} else {
			if typefieldPkg == importPath[1:len(importPath)-1] {
				return importPath, typefieldPkg, typefieldName, definition, nil
			}
		}
	}
	// 	if importPath == imprt {
	// 		base := filepath.Base(importPath)
	// 		// [:removes the last `"` from the package name]
	// 		actualpkg := base[:len(base)-1]
	// 		var definition string
	// 		if actualpkg != "" {
	// 			definition = actualpkg + "." + name
	// 		} else {
	// 			definition = name
	// 		}
	// 		return importPath, actualpkg, name, definition, nil
	return "", "", "", "", fmt.Errorf("Could not locate type %q in import %v.", typefieldName, typefieldPkg)
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
