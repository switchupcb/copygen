package cli

import (
	"bytes"
	"flag"
	"fmt"
	"go/printer"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/switchupcb/copygen/cli/config"
	"github.com/switchupcb/copygen/cli/generator"
	"github.com/switchupcb/copygen/cli/matcher"
	"github.com/switchupcb/copygen/cli/parser"
)

// Environment represents the copygen environment.
type Environment struct {
	YMLPath string // The .yml file path used as a configuration file.
	Output  bool   // Whether to print the generated code to stdout.
}

// CLI runs the copygen command and returns its exit status.
func CLI() int {
	var env Environment

	if err := env.parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}

	if err := env.run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}

// parseArgs parses the provided command line arguments.
func (e *Environment) parseArgs() error {
	// define the command line arguments.
	var (
		ymlpath = flag.String("yml", "", "The path to the .yml flag used for code generation (from the current working directory).")
		output  = flag.Bool("o", false, "Use -o to print generated code to the screen.")
	)

	// parse the command line arguments.
	flag.Parse()

	if !strings.HasSuffix(*ymlpath, ".yml") {
		return fmt.Errorf("you must specify a .yml configuration file using -yml")
	}

	e.YMLPath = *ymlpath
	e.Output = *output

	return nil
}

func (e *Environment) run() error {
	// The configuration file is loaded (.yml)
	gen, err := config.LoadYML(e.YMLPath)
	if err != nil {
		return err
	}

	// The data file is parsed (.go)
	if err = parser.Parse(gen); err != nil {
		return fmt.Errorf("%w", err)
	}

	// The matcher is run on the parsed data (to create the objects used during generation).
	if err = matcher.Match(gen); err != nil {
		return fmt.Errorf("%w", err)
	}

	// Check for used imports.
	usedImports := map[string]bool{}
	for _, function := range gen.Functions {
		for _, fromType := range function.To {
			usedImports[fromType.Field.Import] = true
		}
	}

	// Add new imports if needed
	for path, name := range gen.ImportsByPath {
		if !gen.AlreadyImported[path] && usedImports[path] {
			astutil.AddNamedImport(gen.Fileset, gen.SetupFile, name, path)
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := printer.Fprint(buf, gen.Fileset, gen.SetupFile); err != nil {
		return fmt.Errorf("an error occurred writing the code that will be kept after generation\n%v", err)
	}

	gen.Keep = buf.Bytes()

	// The generator is used to generate code.
	if err = generator.Generate(gen, e.Output); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
