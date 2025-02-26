// Package copygen contains the setup information for copygen generated code.
package copygen

import (
	"github.com/switchupcb/copygen/examples/map/domain"
	"github.com/switchupcb/copygen/examples/map/models"
)

// Copygen defines the functions that are generated.
type Copygen interface {
	A(*models.Account)
	B(*models.User)
	// custom comment
	C(*domain.Account)
	// type basic
	D(int)
	// type basic
	E(string)
	// type basic
	G(float64)
	// type alias
	F(byte)
	H(rune)
	//
	Z()
}
