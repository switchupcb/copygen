package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// parsedFieldName represents the identification data of a parsed *ast.Field.
type parsedFieldName struct {
	pkg  string
	name string
	ptr  string
}

// astParseFieldName parses an *ast.Field (node) for its package, name, and pointer value.
func astParseFieldName(field ast.Node) parsedFieldName {
	var result parsedFieldName

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

	return result
}

// TypeDeclaration represents the information related to a type's declaration.
type TypeDeclaration struct {
	// The package of the type declaration.
	Package *packages.Package

	// The file the type declaration is located in.
	File *ast.File

	// The *ast.TypeSpec of the type declaration.
	TypeSpec *ast.TypeSpec
}

// Locater represents the parameters for locating a Type Declaration.
type Locater struct {
	SetupFile  *ast.File
	Package    string
	Definition string
}

// astLocateTypeDecl uses a setup file (and its imports) to locate the declaration of a type (package.name).
func (p *Parser) astLocateTypeDecl(ltr *Locater) (*TypeDeclaration, error) {
	var ts *ast.TypeSpec

	// check the setup file.
	if ltr.SetupFile.Name.Name == ltr.Package {
		ts, _ = astTypeSearch(ltr.SetupFile, ltr.Definition)
		if ts != nil {
			return &TypeDeclaration{
				Package:  p.LastLocated,
				File:     ltr.SetupFile,
				TypeSpec: ts,
			}, nil
		}
	}

	// check the imports of the setup file.
	for _, setImport := range ltr.SetupFile.Imports {
		// use the exact import (by skipping non-matches) in the case of an alias (i.e `c` in `c "strconv"`)
		if setImport.Name != nil && setImport.Name.Name != ltr.Package {
			continue
		}

		// load a package in the setup file
		pkgs, err := p.loadPackage(setImport.Path.Value[1 : len(setImport.Path.Value)-1])
		if err != nil {
			return nil, err
		}

		// search through the package for a type
		for _, pkg := range pkgs {
			// use the exact package (by skipping non-matches) in the case of no alias
			if pkg.Name != ltr.Package {
				continue
			}

			for _, astFile := range pkg.Syntax {
				ts, _ = astTypeSearch(astFile, ltr.Definition)
				if ts != nil {
					p.LastLocated = pkg

					return &TypeDeclaration{
						Package:  pkg,
						File:     astFile,
						TypeSpec: ts,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("the type declaration for the Field %q of package %q could not be found in the AST.\nIs the imported package up to date?", ltr.Definition, ltr.Package)
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

	return nil, fmt.Errorf("ast: the type %q could not be found in the Abstract Syntax Tree", typename)
}

// parsedDefinition represents the result of a parsed definition.
type parsedDefinition struct {
	err      error
	imprt    string
	pkg      string
	typename string
}

// parseDefinition determines the actual import, package, and name of a field based on its *types.Var definition.
func (p *Parser) parseDefinition(definition string) parsedDefinition {
	var pd parsedDefinition

	// remove pointers
	definition = strings.TrimPrefix(definition, "*")
	splitdefinition := strings.Split(definition, ".")

	// determine the import
	pd.imprt = strings.Join(splitdefinition[:len(splitdefinition)-1], ".")

	// determine the package
	// (i.e `log` in `log.Logger`, `models` in `github.com/.../models.Account`, `models` in `*github.com/.../models/v1.Example`)
	pkgs, err := p.loadPackage(pd.imprt)
	if err != nil {
		pd.err = err

		return pd
	}

	for _, pkg := range pkgs {
		pd.pkg = pkg.Name
	}

	if pd.pkg == "" {
		pd.err = fmt.Errorf("an error occurred determining the package of definition %q", definition)

		return pd
	}

	// determine the typename
	// (i.e `Logger` in `log.Logger`, `DomainUser`)
	base := filepath.Base(definition)
	splitbase := strings.Split(base, ".")

	if len(splitbase) == 1 {
		pd.typename = base
	} else {
		pd.typename = splitbase[1]
	}

	return pd
}
