package generator

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/models"
)

// Generate creates the file with generated code (with gofmt).
func Generate(gen *models.Generator, output bool) error {
	// generate code
	var content string
	header, err := interpreter.Header(gen)
	if err != nil {
		return fmt.Errorf("An error occurred while generating the header.\n%v", err)
	}
	content += header + "\n"

	function, err := interpreter.Function(gen)
	if err != nil {
		return fmt.Errorf("An error occurred while generating a function.\n%v", err)
	}
	content += function + "\n"
	if output {
		fmt.Println(content)
	}

	// gofmt
	data := []byte(content)
	fmtcontent, err := format.Source(data)
	if err != nil {
		return fmt.Errorf("An error occurred while formatting the generated code.\n%v\nUse -o to view output.", err)
	}

	// determine actual filepath
	absfilepath, err := filepath.Abs(gen.Loadpath)
	if err != nil {
		return fmt.Errorf("An error occurred while determing the absolute file path of the generated file.\n%v", absfilepath)
	}
	absfilepath = path.Join(filepath.Dir(absfilepath), gen.Filepath)

	// create file
	if err := os.WriteFile(absfilepath, fmtcontent, 0222); err != nil {
		return fmt.Errorf("An error occurred creating the file.\n%v", err)
	}
	return nil
}
