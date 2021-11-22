// Package parser parses a setup file's functions, types, and fields using an Abstract Syntax Tree.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"

	"github.com/switchupcb/copygen/cli/models"
)

// Parser represents a parser that parses Abstract Syntax Tree data into models.
type Parser struct {
	ImportsByName map[string]string // Map of imports to its alias.
	ImportsByPath map[string]string // Map of imports to its alias.

	// The parser options contain options located in the entire setup file.
	Options OptionMap

	// The fileset of the parser.
	Fileset *token.FileSet

	// The setup file as an Abstract Syntax Tree.
	SetupFile *ast.File

	// Path to root of go project
	goProjectPath string

	// The ast.Node of the `type Copygen Interface`.
	Copygen *ast.InterfaceType

	// A key value cache used to reduce the amount of package load operations during a field search.
	pkg *packages.Package

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
	p := Parser{Setpath: absfilepath, goProjectPath: getProjectPath(absfilepath)}

	config := &packages.Config{
		Mode: pLoadMode,
		Dir:  p.goProjectPath,
	}

	pkgs, err := packages.Load(config, "file="+p.Setpath)
	if err != nil {
		return fmt.Errorf("the setup file's package could not be loaded correctly: %v, %w", p.Setpath, err)
	}
	p.pkg = pkgs[0]

	p.Fileset = token.NewFileSet()

	p.SetupFile, err = parser.ParseFile(p.Fileset, absfilepath, nil, parser.ParseComments)
	imports := astutil.Imports(p.pkg.Fset, p.SetupFile)

	p.ImportsByName = map[string]string{}
	p.ImportsByPath = map[string]string{}
	alreadyImported := map[string]bool{}
	for i := range imports {
		for _, imp := range imports[i] {
			ipath := imp.Path.Value[1 : len(imp.Path.Value)-1]
			name := path.Base(ipath)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			alreadyImported[ipath] = true
			p.ImportsByName[name] = ipath
			p.ImportsByPath[ipath] = name
		}
	}

	if err != nil {
		return fmt.Errorf("an error occurred parsing the specified .go setup file: %v\n%v", gen.Setpath, err)
	}

	p.Options = make(OptionMap)
	p.fieldcache = make(map[string]*models.Field)

	if p.SetupFile.Name == nil {
		return fmt.Errorf("the setup file must declare a package: %v", p.Setpath)
	}

	// Traverse the Abstract Syntax Tree.
	err = p.Traverse(gen)
	if err != nil {
		return err
	}

	gen.Fileset = p.Fileset
	gen.SetupFile = p.SetupFile

	// Add new imports if needed
	for path, name := range p.ImportsByPath {
		if !alreadyImported[path] {
			astutil.AddNamedImport(gen.Fileset, gen.SetupFile, name, path)
		}
	}
	return nil
}

// pLoadMode represents the load mode required for sufficient information during package load.
const pLoadMode = packages.NeedName + packages.NeedImports + packages.NeedTypes + packages.NeedSyntax + packages.NeedTypesInfo

// getProjectPath go up on folders to get directory with `go.mod` file.
func getProjectPath(ppath string) string {
	ppath = strings.Replace(ppath, "\\", "/", -1)
	for ; ppath != ""; ppath = path.Dir(ppath) {
		if _, err := os.Stat(path.Join(ppath, "go.mod")); err == nil {
			break
		}
	}
	return ppath
}
