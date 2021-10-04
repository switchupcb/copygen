// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID    int
	Name  string
	Email string
	User  DomainUser // The fields of DomainUser operate at depth level 1.
}

// DomainUser represents a user in relation to the business logic.
type DomainUser struct {
	UserID   int
	Username string
}
