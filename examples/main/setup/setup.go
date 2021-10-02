// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	c "strconv"

	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// custom: see table in the README for options
	ModelsToDomain(models.Account, models.User) *domain.Account
}

/* Define the fields this converter is applied to using regex. If unspecified, converters are applied to all valid fields. */
// convert: models.User.ID
// Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return c.Itoa(i)
}
