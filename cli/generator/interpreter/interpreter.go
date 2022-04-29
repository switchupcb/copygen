// Package interpreter interprets template code at runtime.
package interpreter

import (
	"fmt"
	"go/build"
	"os"
	"reflect"

	"github.com/switchupcb/copygen/cli/generator/interpreter/extract"
	"github.com/switchupcb/yaegi/interp"
	"github.com/switchupcb/yaegi/stdlib"
)

// InterpretFunction loads a template symbol from an interpreter.
func InterpretFunction(filepath, symbol string) (*reflect.Value, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred loading a template file: %v\nIs the relative or absolute filepath set correctly?\n%w", filepath, err)
	}

	// setup the interpreter
	goCache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("an error occurred loading the template file. Is the GOCACHE set in `go env`?\n%w", err)
	}

	i := interp.New(interp.Options{GoPath: os.Getenv("GOPATH"), GoCache: goCache, GoToolDir: build.ToolDir})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("an error occurred loading the template stdlib libraries\n%w", err)
	}

	// models.types created by the compiled binary are different from models.types created by the interpreter at runtime.
	// pass the compiled models.types to the interpreter
	if err := i.Use(extract.Symbols); err != nil {
		return nil, fmt.Errorf("an error occurred loading the template models library\n%w", err)
	}

	// load the source
	if _, err := i.Eval(string(file)); err != nil {
		return nil, fmt.Errorf("an error occurred evaluating the template file\n%w", err)
	}

	// load the func from the interpreter
	v, err := i.Eval(symbol)
	if err != nil {
		return nil, fmt.Errorf("an error occurred evaluating a template function. Is it located in the file?\n%w", err)
	}

	return &v, nil
}
