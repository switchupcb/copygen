// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	c "strconv"

	"github.com/switchupcb/copygen/examples/error/domain"
	"github.com/switchupcb/copygen/examples/error/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	// custom see table in the README for options
	ModelsToDomain(*models.Account, *models.User) *domain.Account
}

/* Define the function and field this converter is applied to using regex. */
// convert .* models.User.UserID
// Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return c.Itoa(i)
}
