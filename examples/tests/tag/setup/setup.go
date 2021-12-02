// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/tests/tag/domain"
	"github.com/switchupcb/copygen/examples/tests/tag/models"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// depth domain.Account 2
	// depth models.User 1
	// tag domain.Account api
	ModelsToDomain(models.Account, models.User) *domain.Account
}
