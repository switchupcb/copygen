// Package parser parses a setup file's functions, types, and fields using an Abstract Syntax Tree.
package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
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
	p.SetupFile, err = parser.ParseFile(p.Fileset, absfilepath, nil, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("An error occurred parsing the specified .go setup file: %v.\n%v", gen.Setpath, err)
	}

	gen.Imports = p.parseImports()
	gen.Functions, err = p.parseFunctions()
	if err != nil {
		return err
	}
	gen.Keep, err = p.parseKeep(p.Fileset, p.SetupFile)
	if err != nil {
		return err
	}
	return nil
}

// parseImports parses the AST for imports in the setup file.
func (p *Parser) parseImports() map[string]string {
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
	return p.Imports
}

// parseKeep parses the generator's setup file for data that is kept in the generated file.
// TODO: Implement Keep
func (p *Parser) parseKeep(fileset *token.FileSet, file *ast.File) (string, error) {
	var keep []byte
	buffer := bytes.NewBuffer(keep)
	ast.FilterFile(file, func(s string) bool {
		if s == "Copygen" {
			return false
		}
		// Keep all types that are not Copygen.
		return true
	})

	if err := printer.Fprint(buffer, fileset, file); err != nil {
		return "", err
	}
	return string(buffer.Bytes()), nil
}
