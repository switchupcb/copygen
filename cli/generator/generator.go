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
const writeFileMode = 0222

// Generate outputs the generated code (with gofmt).
func Generate(gen *models.Generator, output bool, write bool) (string, error) {
	content, err := generateCode(gen)
	if err != nil {
		return "", fmt.Errorf("an error occurred while generating code\n%w", err)
	}

	// gofmt
	data := []byte(content)
	fmtcontent, err := format.Source(data)
	if err != nil {
		if output {
			fmt.Println(content)
			return content, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
		}

		return content, fmt.Errorf("an error occurred while formatting the generated code.\n%w\nUse -o to view output", err)
	}

	code := string(fmtcontent)
	if output {
		fmt.Println(code)
		return code, nil
	}

	if write {
		// determine actual filepath
		absfilepath, err := filepath.Abs(gen.Loadpath)
		if err != nil {
			return code, fmt.Errorf("an error occurred while determining the absolute file path of the generated file\n%v", absfilepath)
		}

		absfilepath = filepath.Join(filepath.Dir(absfilepath), gen.Outpath)

		// create file
		if err := os.WriteFile(absfilepath, fmtcontent, writeFileMode); err != nil {
			return code, fmt.Errorf("an error occurred creating the file.\n%w", err)
		}
	}

	return code, nil
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
		return "", fmt.Errorf("an error occurred loading the absolute filepath of template path %v from the cwd %v\n%w", gen.Loadpath, gen.Tempath, err)
	}

	v, err := interpreter.InterpretFunction(abstempath, GenerateFunction)
	if err != nil {
		return "", fmt.Errorf("%w", err)
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
