// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/_tests/cyclic/domain"
	"github.com/switchupcb/copygen/examples/_tests/cyclic/duplicate"
	"github.com/switchupcb/copygen/examples/_tests/cyclic/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	ModelsToDomain(*models.Account, *models.User) *domain.Account
	DuplicateCyclic(*duplicate.Account) *domain.Account
	SuperCyclic(domain.CyclicInterface) *domain.CyclicInterface
}
