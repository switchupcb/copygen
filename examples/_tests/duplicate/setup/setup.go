// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/_tests/duplicate/domain"
	"github.com/switchupcb/copygen/examples/_tests/duplicate/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	ModelsToDomain(models.Account, models.User) *domain.Account
}
