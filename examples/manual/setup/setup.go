// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	c "strconv"

	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// map models.Account.ID *domain.Account.ID
	// map models.Account.Name *domain.Account.Name
	// map models.User.UserID *domain.Account.UserID
	ModelsToDomain(models.Account, models.User) *domain.Account
}

/* Define the function and field this converter is applied to using regex. */
// convert .* models.User.UserID
// Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return c.Itoa(i)
}
