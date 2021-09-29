// Package yml loads setup information from an external file.
package yml

import (
	"fmt"
	"os"

	goloader "github.com/switchupcb/copygen/cli/loader/go"
	"github.com/switchupcb/copygen/cli/models"
	"gopkg.in/yaml.v3"
)

// Parser represents a YML parser that loads properties into the program models.
type Parser struct {
	YML       YML              // The YML that is parsed.
	AST       goloader.AST     // The Abstract Syntax Tree object used during matching.
	Generator models.Generator // The generator that information is parsed to.

}

// LoadYML loads a .yml file into a Generator.
func LoadYML(filepath string) (*models.Generator, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("The specified .yml filepath doesn't exist: %v.\n%v", filepath, err)
	}

	var p Parser
	err = yaml.Unmarshal(file, &p.YML)
	if err != nil {
		return nil, fmt.Errorf("There is an issue with the provided .yml file: %v\n%v", filepath, err)
	}

	gen := p.ParseYML()
	if err != nil {
		return nil, err
	}
	gen.Loadpath = filepath
	return gen, nil
}

// ParseYML parses a YML into a Generator.
func (p *Parser) ParseYML() *models.Generator {
	p.Generator.Setpath = p.YML.Generated.Setup
	p.Generator.Outpath = p.YML.Generated.Output
	p.Generator.Package = p.YML.Generated.Package
	p.Generator.Template.Funcpath = p.YML.Templates.Function
	p.Generator.Template.Headpath = p.YML.Templates.Header
	return &p.Generator
}
