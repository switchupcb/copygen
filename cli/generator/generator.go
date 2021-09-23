// generator generates code based on a template.
package generator

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/models"
)

// Generate creates the file with generated code (with gofmt).
func Generate(g *models.Generator) error {
	// generate code
	var file string
	file += Header(g) + "\n"
	for _, function := range g.Functions {

		file += Function(&function)
		file += "\n"
	}

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
