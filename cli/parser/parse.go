// Package parser parses a setup file's functions, types, and fields using an Abstract Syntax Tree.
package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
	"golang.org/x/tools/go/packages"
)

// Parser represents a parser that parses Abstract Syntax Tree data into models.
type Parser struct {
	Config  Config
	Options Options
	Pkgs    []*packages.Package
}

// Config represents a Parser's configuration.
type Config struct {
	// SetupFile represents the setup file as an Abstract Syntax Tree.
	SetupFile *ast.File

	// SetupPkg represent the setup file's package.
	SetupPkg *packages.Package

	// Fileset represents the parser's fileset.
	Fileset *token.FileSet
}

// Options represents a parser's options.
type Options struct {
	// commentOptionMap represents a map of comments (as text) to an option.
	CommentOptionMap map[string]*options.Option

	// convertOptions represents a global list of convert options (for convert functions).
	ConvertOptions []*options.Option
}

// GLOBAL VARIABLES.
var (
	// fieldcache represents a map of `go/types` Type strings to models.Field.
	//
	// fieldcache is used to prevent cyclic fields from incorrect assignment.
	//
	// fieldcache improves performance by parsing a unique type definition once per runtime.
	// definitions remain constant UNLESS the user modifies their modules during runtime.
	fieldcache map[string]*models.Field

	// setupPkgPath represents the current path of the setup file's package.
	//
	// setupPkgPath is used to remove package references from types that will be
	// used in the generated file's package (equal to the setup file's package).
	//
	// setupPkgPath is referenced while parsing collected type definitions for collection fields,
	// and while setting package references for non-collection fields after parsing.
	//
	// i.e `Collections` parsed as `copygen.Collections` in the setup file's package copygen,
	// output as `Collections` in the generated file's package copygen.
	setupPkgPath string

	// outputPkgPath represents the generated file's package path.
	//
	// outputPkgPath is used to remove package references from types that are imported
	// (in the setup file) from the generated file's package.
	//
	// outputPkgPath is referenced while parsing collected type definitions for collection fields,
	// and while setting package references for non-collection fields after parsing.
	outputPkgPath string

	// aliasImportMap represents a map of import paths to package names.
	//
	// aliasImportMap is used to assign the correct package reference to an aliased field.
	//
	// aliasImportMap is referenced while parsing collected type definitions for collection fields,
	// and while setting package references for non-collection fields after parsing.
	aliasImportMap map[string]string
)

// SetupCache sets up the parser's global cache.
func SetupCache() {
	if fieldcache == nil {
		fieldcache = make(map[string]*models.Field)
	}
}

// ResetCache resets the parser's global cache.
func ResetCache() {
	fieldcache = make(map[string]*models.Field)
}

// parserLoadMode represents the load mode required for sufficient information during package load.
const parserLoadMode = packages.NeedName + packages.NeedImports + packages.NeedDeps + packages.NeedTypes + packages.NeedSyntax + packages.NeedTypesInfo

// Parse parses a generator's setup file.
func Parse(gen *models.Generator) error {
	var err error
	p := new(Parser)
	p.Config.Fileset = token.NewFileSet()
	p.Config.SetupFile, err = parser.ParseFile(p.Config.Fileset, gen.Setpath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("an error occurred parsing the specified .go setup file: %v\n%w", gen.Setpath, err)
	}

	// Parse the setup file's `type Copygen Interface` for the Keep (and create Options in the process).
	if err := p.Keep(p.Config.SetupFile); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Analyze a new `type Copygen Interface` to create models.Function and models.Field objects.
	cfg := &packages.Config{Mode: parserLoadMode}
	p.Pkgs, err = packages.Load(cfg, "file="+gen.Setpath)
	if err != nil {
		return fmt.Errorf("an error occurred while loading the packages for types.\n%w", err)
	}
	p.Config.SetupPkg = p.Pkgs[0]

	// determine the output file package path.
	outputPkgs, _ := packages.Load(&packages.Config{Mode: packages.NeedName}, "file="+gen.Outpath)
	if len(outputPkgs) > 0 {
		outputPkgPath = outputPkgs[0].PkgPath
	}

	// set the aliasImportMap.
	aliasImportMap = make(map[string]string, len(p.Config.SetupFile.Imports))
	for _, imp := range p.Config.SetupFile.Imports {
		if imp.Name != nil {
			aliasImportMap[imp.Path.Value[1:len(imp.Path.Value)-1]] = imp.Name.Name
		}
	}

	// find a new instance of a `type Copygen interface` AST from the setup file's
	// loaded go/types package (containing different *ast.Files from the Keep)
	// since the parsed `type Copygen interface` has its comments removed.
	var newCopygen *ast.InterfaceType
	for _, decl := range p.Config.SetupPkg.Syntax[0].Decls {
		switch declaration := decl.(type) {
		case *ast.GenDecl:
			if it, ok := assertCopygenInterface(declaration); ok {
				newCopygen = it
				break
			}
		}
	}

	if newCopygen == nil {
		return fmt.Errorf("the \"type Copygen interface\" could not be found in the setup file")
	}

	// create models.Function objects.
	SetupCache()
	if gen.Functions, err = p.parseFunctions(newCopygen); err != nil {
		return fmt.Errorf("%w", err)
	}

	// rename non-collection fields' packages using imports.
	setPackages(gen)

	// Write the Keep.
	buf := new(bytes.Buffer)
	buf.WriteString("// Code generated by github.com/switchupcb/copygen\n// DO NOT EDIT.\n\n")
	if err := printer.Fprint(buf, p.Config.Fileset, p.Config.SetupFile); err != nil {
		return fmt.Errorf("an error occurred writing the code that will be kept after generation\n%w", err)
	}
	gen.Keep = buf.Bytes()

	// reset global variables.
	setupPkgPath = ""
	outputPkgPath = ""
	aliasImportMap = nil

	return nil
}

// setPackages sets the packages for all fields in a generator using names from the setup file.
func setPackages(gen *models.Generator) {
	for _, function := range gen.Functions {
		functionTypes := [][]models.Type{
			function.From,
			function.To,
		}

		for _, types := range functionTypes {
			for _, t := range types {
				for _, field := range t.Field.AllFields(nil, nil) {

					// a generated file's package == setup file's package.
					//
					// when the field is defined in the setup file (i.e `Collection`),
					// it will be parsed with the setup file's package (i.e `copygen.Collection`).
					//
					// do NOT reference it by package in the generated file (i.e `Collection`).
					if field.Import == setupPkgPath {
						field.Package = ""
						continue
					}

					// when a setup file imports the package it will output to,
					// do NOT reference the fields defined in the output package, by package.
					if outputPkgPath != "" && field.Import == outputPkgPath {
						field.Package = ""
						continue
					}

					// when a field's import uses an alias, reassign the package reference.
					if aliasPkg, ok := aliasImportMap[field.Import]; ok {
						field.Package = aliasPkg
						continue
					}
				}
			}
		}
	}
}
