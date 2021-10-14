package generator

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/generator/interpreter"
	"github.com/switchupcb/copygen/cli/models"
)

// Generate creates the file with generated code (with gofmt).
func Generate(gen *models.Generator, output bool) error {
	// generate code
	content, err := interpreter.Generate(gen)
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
