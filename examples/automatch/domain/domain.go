// The domain package contains business logic models.
package domain

import "log"

type T int

// Account represents a user account.
type Account struct {
	ID    int
	Name  string
	Email string
	Basic T
	User  DomainUser
	Depth log.Logger
}

// DomainUser represents a domain user.
type DomainUser struct {
	UserID   string
	Name     string
	UserData map[string]interface{}
}
