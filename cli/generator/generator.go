package generator

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/traefik/yaegi/interp"
)

// Generate creates the file with generated code (with gofmt).
func Generate(g *models.Generator) error {
	// generate code
	var file string
	header, err := Header(g)
	if err != nil {
		return err
	}
	file += header + "\n"

	function, err := Function(g)
	if err != nil {
		return err
	}
	file += function + "\n"

	// gofmt
	data := []byte(file)
	content, err := format.Source(data)
	if err != nil {
		return err
	}

	// determine actual filepath
	fpath, err := filepath.Abs(g.Loadpath)
	if err != nil {
		return err
	}
	fpath = path.Join(filepath.Dir(fpath), g.Filepath)

	// create file
	if err := os.WriteFile(fpath, content, 0222); err != nil {
		return fmt.Errorf("There was an error creating the file.\n%v", err)
	}
	return nil
}

// Interpret loads a template function using an interpreter.
func Interpret(lpath string, tpath, symbol string) (func(interface{}) string, error) {
	// determine actual filepath
	fpath, err := filepath.Abs(lpath)
	if err != nil {
		return nil, err
	}
	fpath = path.Join(filepath.Dir(fpath), tpath)

	// read the file
	file, err := os.ReadFile(fpath)
	if err != nil {
		return nil, fmt.Errorf("The specified template file doesn't exist: %v\n%v", fpath, err)
	}
	source := string(file)

	// interpret the source file
	i := interp.New(interp.Options{})
	if _, err := i.Eval(source); err != nil {
		return nil, fmt.Errorf("There was an error running the template file: %v\n%v", fpath, err)
	}

	// get the required symbol (package.function)
	v, err := i.Eval(symbol)
	if err != nil {
		return nil, fmt.Errorf("There was an error loading a template function.\n%v", err)
	}
	return v.Interface().(func(interface{}) string), nil
}
