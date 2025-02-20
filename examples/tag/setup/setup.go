// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/tag/domain"
	"github.com/switchupcb/copygen/examples/tag/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	// tag .* api
	ModelsToDomain(*models.Account, *models.User) *domain.Account
}
