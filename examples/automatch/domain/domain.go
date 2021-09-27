// The domain package contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID    int
	Name  string
	Email string
	User  DomainUser // The fields of DomainUser operate at depth level 1.
}

type DomainUser struct {
	UserID int
	Name   string
}
