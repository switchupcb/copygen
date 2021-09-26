package interpreter

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/traefik/yaegi/interp"
)

// InterpretFunc loads a template package.function using an interpreter.
func InterpretFunc(lpath string, tpath string, symbol string) (func(interface{}) string, error) {
	i, err := interpretFile(lpath, tpath, symbol)
	if err != nil {
		return nil, err
	}

	v, err := i.Eval(symbol)
	if err != nil {
		return nil, fmt.Errorf("There was an error loading a template function.\n%v", err)
	}
	return v.Interface().(func(interface{}) string), nil
}

// interpretFile loads a template file using an interpreter.
func interpretFile(lpath string, tpath, symbol string) (*interp.Interpreter, error) {
	// determine actual filepath
	fpath, err := filepath.Abs(lpath)
	if err != nil {
		return nil, err
	}
	fpath = path.Join(filepath.Dir(fpath), tpath)

	// read the file
	file, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("The specified template file for the template function %v doesn't exist: %v\n", symbol, fpath)
	}
	source := string(file)

	// interpret the source file
	i := interp.New(interp.Options{})
	if _, err := i.Eval(source); err != nil {
		return nil, fmt.Errorf("There was an error running the template file: %v\n%v", fpath, err)
	}
	return i, nil
}
