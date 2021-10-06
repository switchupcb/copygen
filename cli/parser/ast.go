package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

type parsedASTFieldName struct {
	pkg  string
	name string
	def  string
	ptr  string
}

// parseASTFieldName parses an *ast.Field (node) for its package, name, definition, and pointer value.
func parseASTFieldName(field ast.Node) parsedASTFieldName {
	var result parsedASTFieldName

	ast.Inspect(field, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.SelectorExpr:
			// FieldInfo is always in a selector expression.
			result.pkg += x.X.(*ast.Ident).Name // 'log' in 'Field log.Logger'
			result.name += x.Sel.Name           // 'Logger' in 'Field log.Logger'
			return false
		case *ast.StarExpr:
			result.ptr += "*"
			return true
		default:
			return true
		}
	})
	if result.pkg != "" {
		result.def = fmt.Sprintf("%s.%s", result.pkg, result.name)
	} else {
		result.def = result.name
	}
	return result
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
	return "", fmt.Errorf("could not locate type %q in file import %v", pkg+" "+name, fileImport)
}

// astTypeSearch searches through an ast.File for ast.Types.
func astTypeSearch(file *ast.File, typename string) (*ast.TypeSpec, error) {
	for _, decl := range file.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

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
	return nil, fmt.Errorf("the type %q could not be found in the Abstract Syntax Tree", typename)
}
