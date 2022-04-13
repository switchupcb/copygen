// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID    int
	Name  string
	Email string
	User  DomainUser // The fields of DomainUser operate at depth level 1.
}

// DomainUser represents a user in relation to business logic.
type DomainUser struct {
	UserID   int
	Username string
	Password Password // The fields of Password operate at depth level 2.
}

// Password represents a password in relation to business logic.
type Password struct {
	Password string
	Hash     string
	Salt     string
}
