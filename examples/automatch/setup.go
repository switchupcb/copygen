// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// depth domain.Account.*(?!\.) 1
	// depth models.User.*(?!\.) 2
	ModelsToDomain(models.Account, models.User) domain.Account
}
