package interpreter

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

//nolint:unused // interpretFunc loads a template package.function into an interpreter.
func interpretFunc(loadpath string, templatepath, symbol string) (*reflect.Value, error) {
	// determine actual filepath
	absfilepath, err := filepath.Abs(loadpath)
	if err != nil {
		return nil, err
	}
	absfilepath = path.Join(filepath.Dir(absfilepath), templatepath)

	// read the file
	file, err := os.ReadFile(absfilepath)
	if err != nil {
		return nil, fmt.Errorf("The specified template file for the template function %q doesn't exist: %v\nIs the relative or absoute filepath set correctly?", symbol, absfilepath)
	}
	source := string(file)

	// setup the interpreter
	/* goCache, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("An error occurred loading the template file. Is the GOCACHE set in `go env`?", err)
	} */

	// create the interpreter
	i := interp.New(interp.Options{GoPath: os.Getenv("GOPATH") /*, GoCache: goCache, GoToolDir: build.ToolDir*/})
	if err := i.Use(stdlib.Symbols); err != nil {
		return nil, fmt.Errorf("An error occurred loading the template libraries.\n%v", err)
	}
	if _, err := i.Eval(source); err != nil {
		return nil, fmt.Errorf("An error occurred loading the template file: %v\n%v", absfilepath, err)
	}

	// get the func from the interpreter
	v, err := i.Eval(symbol)
	if err != nil {
		return nil, fmt.Errorf("An error occurred loading a template function.\n%v", err)
	}
	return &v, nil
}
