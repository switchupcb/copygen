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
}

// Runs the copygen command and returns its exit status.
func Run() int {
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
	gen, err := config.LoadYML(e.YMLPath)
	if err != nil {
		return err
	}

	if err = parser.Parse(gen); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = matcher.Match(gen); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = generator.Generate(gen, e.Output); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
