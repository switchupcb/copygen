package interpreter

import (
	"fmt"
	"go/build"
	"os"
	"reflect"

	"github.com/switchupcb/copygen/cli/generator/interpreter/extract"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// InterpretFunction loads a template symbol from an interpreter.
func InterpretFunction(filepath, symbol string) (*reflect.Value, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("an error occurred loading a template file: %v\nIs the relative or absoute filepath set correctly?\n%v", filepath, err)
	}

	// setup the interpreter
	goCache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("an error occurred loading the template file. Is the GOCACHE set in `go env`?\n%v", err)
	}

	i := interp.New(interp.Options{GoPath: os.Getenv("GOPATH"), GoCache: goCache, GoToolDir: build.ToolDir})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("an error occurred loading the template stdlib libraries\n%v", err)
	}

	// models.types created by the compiled binary are different from models.types created by the interpreter at runtime.
	// pass the compiled models.types to the interpreter
	if err := i.Use(extract.Symbols); err != nil {
		return nil, fmt.Errorf("an error occurred loading the template models library\n%v", err)
	}

	// load the source
	if _, err := i.Eval(string(file)); err != nil {
		return nil, fmt.Errorf("an error occurred evaluating the template file\n%v", err)
	}

	// load the func from the interpreter
	v, err := i.Eval(symbol)
	if err != nil {
		return nil, fmt.Errorf("an error occurred evaluating a template function. Is it located in the file?\n%v", err)
	}

	return &v, nil
}
