package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/switchupcb/copygen/cli/config"
	"github.com/switchupcb/copygen/cli/generator"
	"github.com/switchupcb/copygen/cli/matcher"
	"github.com/switchupcb/copygen/cli/parser"
)

// Environment represents the copygen environment.
type Environment struct {
	YMLPath string // The .yml file path used as a configuration file.
	Output  bool   // Whether to print the generated code to stdout.
	Write   bool   // Whether to write the generated code to a file.
}

// CLI runs copygen from a Command Line Interface and returns the exit status.
func CLI() int {
	var env Environment

	if err := env.parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}

	if _, err := env.Run(); err != nil {
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
	e.Write = true

	return nil
}

// Run runs copygen programmatically using the given Environment's YMLPath.
func (e *Environment) Run() (string, error) {
	// The configuration file is loaded (.yml)
	gen, err := config.LoadYML(e.YMLPath)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// The data file is parsed (.go)
	if err = parser.Parse(gen); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// The matcher is run on the parsed data (to create the objects used during generation).
	if !gen.Options.Matcher.Skip {
		if err = matcher.Match(gen); err != nil {
			return "", fmt.Errorf("%w", err)
		}
	}

	// The generator is used to generate code.
	code, err := generator.Generate(gen, e.Output, e.Write)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return code, nil
}
