// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/_tests/cast/domain"
	"github.com/switchupcb/copygen/examples/_tests/cast/models"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// cast domain.ReversedString .*
	// cast string SuperString
	ModelsToDomain(*models.Account) *domain.Account
}
