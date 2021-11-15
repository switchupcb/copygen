// Package parser parses a setup file's functions, types, and fields using an Abstract Syntax Tree.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/switchupcb/copygen/cli/models"
)

// Parser represents a parser that parses Abstract Syntax Tree data into models.
type Parser struct {
	// The parser options contain options located in the entire setup file.
	Options OptionMap

	// The fileset of the parser.
	Fileset *token.FileSet

	// The setup file as an Abstract Syntax Tree.
	SetupFile *ast.File

	// The ast.Node of the `type Copygen Interface`.
	Copygen *ast.InterfaceType

	// A key value cache used to reduce the amount of package load operations during a field search.
	pkgcache map[string][]*packages.Package

	// The last package to be loaded by a Locater.
	LastLocated *packages.Package

	// A key value cache used to reduce the amount of AST operations during a field search.
	fieldcache map[string]*models.Field

	// The setup filepath.
	Setpath string

	// The option-comments parsed in the OptionMap.
	Comments []*ast.Comment
}

// Parse parses a generator's setup file.
func Parse(gen *models.Generator) error {
	// determine the actual filepath of the setup.go file.
	absfilepath, err := filepath.Abs(filepath.Join(filepath.Dir(gen.Loadpath), gen.Setpath))
	if err != nil {
		return err
	}

	// setup the parser
	p := Parser{Setpath: absfilepath}
	p.Fileset = token.NewFileSet()

	p.SetupFile, err = parser.ParseFile(p.Fileset, absfilepath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("an error occurred parsing the specified .go setup file: %v\n%v", gen.Setpath, err)
	}

	p.Options = make(OptionMap)
	p.fieldcache = make(map[string]*models.Field)
	p.pkgcache = make(map[string][]*packages.Package)

	pkgs, err := p.loadPackage("file=" + p.Setpath)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if p.SetupFile.Name == nil {
			return fmt.Errorf("the setup file must declare a package: %v", p.Setpath)
		} else if p.SetupFile.Name.Name == pkg.Name {
			p.LastLocated = pkg
			break
		}
	}

	if p.LastLocated == nil {
		return fmt.Errorf("the setup file's package could not be loaded correctly: %v", p.Setpath)
	}

	// Traverse the Abstract Syntax Tree.
	err = p.Traverse(gen)
	if err != nil {
		return err
	}

	gen.Fileset = p.Fileset
	gen.SetupFile = p.SetupFile

	imports := astutil.Imports(gen.Fileset, gen.SetupFile)

	gen.ImportsByName = map[string]string{}
	gen.ImportsByPath = map[string]string{}
	gen.AlreadyImported = map[string]bool{}
	for i := range imports {
		for _, imp := range imports[i] {
			ipath := imp.Path.Value[1 : len(imp.Path.Value)-1]
			name := path.Base(ipath)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			gen.AlreadyImported[ipath] = true
			gen.ImportsByName[name] = ipath
			gen.ImportsByPath[ipath] = name
		}
	}
	return nil
}

// pLoadMode represents the load mode required for sufficient information during package load.
const pLoadMode = packages.NeedName + packages.NeedImports + packages.NeedTypes + packages.NeedSyntax + packages.NeedTypesInfo

// loadPackage loads a package.
func (p *Parser) loadPackage(importPath string) ([]*packages.Package, error) {
	if pkgs, exists := p.pkgcache[importPath]; exists {
		return pkgs, nil
	}

	config := &packages.Config{Mode: pLoadMode, Logf: nil}

	pkgs, err := packages.Load(config, importPath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred loading a package from the GOPATH with import: %v.\n%v", importPath, err)
	}

	p.pkgcache[importPath] = pkgs

	return pkgs, nil
}
