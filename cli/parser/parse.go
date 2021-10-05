// Package parser parses a setup file's functions, types, and fields using an Abstract Syntax Tree.
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/models"
)

// Parser represents a parser that parses Abstract Syntax Tree data into models.
type Parser struct {
	// The setup filepath.
	Setpath string

	// The setup file as an Abstract Syntax Tree.
	SetupFile *ast.File

	// The fileset of the parser.
	Fileset *token.FileSet

	// The ast.Node of the `type Copygen Interface`.
	Copygen *ast.InterfaceType

	// The parser options contain options located in the entire setup file.
	Options OptionMap

	// The imports discovered in the set up file (map[packagevar]importpath).
	// In the context of the parser, packagevar refers to the the variable used
	// to reference the package (alias) rather the package's actual name.
	Imports map[string]string
}

// Parse parses a generator's setup file.
func Parse(gen *models.Generator) error {
	// determine the actual filepath of the setup.go file.
	absfilepath, err := filepath.Abs(filepath.Join(filepath.Dir(gen.Loadpath), gen.Setpath))
	if err != nil {
		return err
	}

	p := Parser{Setpath: absfilepath}
	p.Fileset = token.NewFileSet()
	p.SetupFile, err = parser.ParseFile(p.Fileset, absfilepath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("An error occurred parsing the specified .go setup file: %v.\n%v", gen.Setpath, err)
	}

	p.Options = make(OptionMap)
	gen.Keep, err = p.parseKeep()
	if err != nil {
		return err
	} else if p.Copygen == nil {
		return fmt.Errorf("The \"type Copygen interface\" could not be found in the setup file.")
	}
	gen.Functions, err = p.parseFunctions()
	if err != nil {
		return err
	}
	return nil
}

// parseImports parses the AST for imports in the setup file.
func (p *Parser) parseImports() {
	if p.Imports == nil {
		p.Imports = make(map[string]string) // map[packagevar]importpath
	}

	for _, imprt := range p.SetupFile.Imports {
		if imprt.Name != nil { // aliased package (i.e c "strconv")
			p.Imports[imprt.Name.Name] = imprt.Path.Value
		} else {
			base := filepath.Base(imprt.Path.Value)
			// [:removes the last `"` from the package name]
			p.Imports[base[:len(base)-1]] = imprt.Path.Value
		}
	}
}
