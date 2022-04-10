// Package generator generates code.
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	tmpl "text/template"

	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/generator/template"
	"github.com/switchupcb/copygen/cli/models"
)

const (
	GenerateFunction = "template.Generate"
	writeFileMode    = 0222
)

// Generate outputs the generated code (with gofmt).
func Generate(gen *models.Generator, output bool, write bool) (string, error) {
	content, err := generate(gen)
	if err != nil {
		return "", fmt.Errorf("an error occurred while generating code.\n%w", err)
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

// generate determines the method of code generation to use,
// then generates the code.
func generate(gen *models.Generator) (string, error) {
	if gen.Tempath != "" {
		var err error
		ext := filepath.Ext(gen.Tempath)
		gen.Tempath, err = filepath.Abs(filepath.Join(filepath.Dir(gen.Loadpath), gen.Tempath))
		if err != nil {
			return "", fmt.Errorf("an error occurred loading the absolute filepath of template path %v from the cwd %v\n%w", gen.Loadpath, gen.Tempath, err)
		}

		// generate code using a .go template.
		if ext == ".go" {
			return generateCode(gen)
		}

		// generate code using a .tmpl template.
		if ext == ".tmpl" {
			return generateTemplate(gen)
		}

		return "", fmt.Errorf("the provided template is not a `.go` or `.tmpl` file: %v", gen.Tempath)
	}

	// generate code using the default template.
	return template.Generate(gen)
}

// generateCode generates code using the default .go template.
func generateCode(gen *models.Generator) (string, error) {
	// use an interpreted function (from a template file).
	v, err := interpreter.InterpretFunction(gen.Tempath, GenerateFunction)
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

var (
	// funcMap represents a funcMap for text/templates.
	funcMap = tmpl.FuncMap{
		"bytesToString": func(b []byte) string { return string(b) },
	}
)

// generateTemplate generates code using a text/template file (.tmpl).
func generateTemplate(gen *models.Generator) (string, error) {
	file, err := os.ReadFile(gen.Tempath)
	if err != nil {
		return "", fmt.Errorf("the specified .tmpl filepath doesn't exist: %v\n%w", gen.Tempath, err)
	}

	t, err := tmpl.New("").Funcs(funcMap).Parse(string(file))
	if err != nil {
		return "", fmt.Errorf("an error occurred parsing the .tmpl template file: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	if err = t.Execute(buf, gen); err != nil {
		return "", fmt.Errorf("an error occurred executing the .tmpl template file: %w", err)
	}

	return buf.String(), nil
}
