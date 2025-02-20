// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	service "github.com/switchupcb/copygen/examples/main/domain"
	data "github.com/switchupcb/copygen/examples/main/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	ModelsToDomain(*data.Account, *data.User) *service.Account
}
