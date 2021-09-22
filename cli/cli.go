// package cli contains the primary logic of the copygen command-line application.
package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/switchupcb/copygen/cli/loader"
)

// Environment represents the copygen environment.
type Environment struct {
	YML string // The yml file path used as a configuration file.
}

// CLI runs the copygen command and returns its exit status.
func CLI(args []string) int {
	var env Environment
	err := env.parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Argument error: %v\n", err)
		return 2
	}

	if err = env.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

// parseArgs parses the provided command line arguments.
func (e *Environment) parseArgs(args []string) error {
	// define the command line arguments.
	ymlPtr := flag.String("yml", "", "The path to the .yml flag used for code generation (from the current working directory).")

	// parse the command line arguments.
	flag.Parse()

	// yml
	ymlLen := len(*ymlPtr)
	if ymlLen == 0 {
		return fmt.Errorf("No .yml configuration file was specified using -yml.")
	} else if ymlLen < 4 || ".yml" != (*ymlPtr)[ymlLen-4:] {
		return fmt.Errorf("The specified file is not a .yml file.")
	}
	e.YML = *ymlPtr
	return nil
}

func (e *Environment) run() error {
	generator, err := loader.LoadYML(e.YML)
	if err != nil {
		return err
	}

	fmt.Println(*generator)
	return nil
}
