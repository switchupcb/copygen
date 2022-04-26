// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID    int
	Name  string
	Email string
	Info  Cyclic   // Info contains a cyclic field.
	Owner *Account // Owner is a cyclic field.
}

// Cyclic represents a cyclic type (that holds an account).
type Cyclic struct {
	UserID   int
	Username string
	Account  *Account
}

// CyclicInterface represents a cyclic interface (that contains a cyclic func).
type CyclicInterface interface {
	CyclicFunc(CyclicInterface) bool
}
