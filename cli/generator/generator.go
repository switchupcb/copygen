package generator

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"

	"github.com/switchupcb/copygen/cli/generator/template"
	"github.com/switchupcb/copygen/cli/models"
)

// Generate creates the file with generated code (with gofmt).
func Generate(g *models.Generator) error {
	// generate code
	var content string
	header, err := template.Header(g)
	if err != nil {
		return err
	}
	content += header + "\n"

	function, err := template.Function(g)
	if err != nil {
		return err
	}
	content += function + "\n"

	// gofmt
	data := []byte(content)
	fmtcontent, err := format.Source(data)
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
	if err := os.WriteFile(fpath, fmtcontent, 0222); err != nil {
		return fmt.Errorf("There was an error creating the file.\n%v", err)
	}
	return nil
}
