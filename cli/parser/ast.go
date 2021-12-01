package parser

import (
	"path/filepath"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// parsedFieldName represents the identification data of a parsed *ast.Field.
type parsedFieldName struct {
	pkg  string
	name string
	ptr  string
}

// parsedDefinition represents the result of a parsed definition.
type parsedDefinition struct {
	err           error
	imprt         string
	typename      string
	containerType models.ContainerType
	pointer       string
}

// parseDefinition determines the actual import, package, and name of a field based on its *types.Var definition.
func (p *Parser) parseDefinition(definition string) parsedDefinition {
	var pd parsedDefinition

	// remove pointers
	if strings.Index(definition, "[]") == 0 {
		definition = strings.TrimPrefix(definition, "[]")
		pd.containerType = models.ContainerTypeSlice
	}
	if strings.Index(definition, "*") == 0 {
		pd.pointer = "*"
	}
	definition = strings.TrimPrefix(definition, "*")
	splitdefinition := strings.Split(definition, ".")

	// determine the import
	pd.imprt = strings.Join(splitdefinition[:len(splitdefinition)-1], ".")

	// determine the typename
	// (i.e `Logger` in `log.Logger`, `DomainUser`)
	base := filepath.Base(definition)
	splitbase := strings.Split(base, ".")

	pd.typename = splitbase[len(splitbase)-1]

	return pd
}
