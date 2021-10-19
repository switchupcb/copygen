package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/generator/template"
	"github.com/switchupcb/copygen/cli/models"
)

const GenerateFunction = "template.Generate"

// Generate creates the file with generated code (with gofmt).
func Generate(gen *models.Generator, output bool) error {
	content, err := generateCode(gen)
	if err != nil {
		return fmt.Errorf("an error occurred while generating code\n%v", err)
	}

	if output {
		fmt.Println(content)
	}

	// gofmt
	data := []byte(content)

	fmtcontent, err := format.Source(data)
	if err != nil {
		return fmt.Errorf("an error occurred while formatting the generated code.\n%v\nUse -o to view output", err)
	}

	// determine actual filepath
	absfilepath, err := filepath.Abs(gen.Loadpath)
	if err != nil {
		return fmt.Errorf("an error occurred while determining the absolute file path of the generated file\n%v", absfilepath)
	}

	absfilepath = filepath.Join(filepath.Dir(absfilepath), gen.Outpath)

	// create file
	if err := os.WriteFile(absfilepath, fmtcontent, 0222); err != nil { //nolint:gofumpt // ignore
		return fmt.Errorf("an error occurred creating the file.\n%v", err)
	}

	return nil
}

// generateCode determines the func to generate function code.
func generateCode(gen *models.Generator) (string, error) {
	if gen.Tempath == "" {
		content, _ := template.Generate(gen)
		return content, nil
	}

	// use an interpreted function (from a template file)
	abstempath, err := filepath.Abs(filepath.Join(filepath.Dir(gen.Loadpath), gen.Tempath))
	if err != nil {
		return "", fmt.Errorf("an error occurred loading the absolute filepath of template path %v from the cwd %v\n%v", gen.Loadpath, gen.Tempath, err)
	}

	v, err := interpreter.InterpretFunction(abstempath, GenerateFunction)
	if err != nil {
		return "", err
	}

	fn, ok := v.Interface().(func(*models.Generator) (string, error))
	if !ok {
		return "", fmt.Errorf("the template function `Generate` could not be type asserted. Is it a func(*models.Generator) (string, error)?")
	}

	content, err := fn(gen)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return content, nil
}
