// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID       string
	Name     string
	Email    string
	Password string
	Other    string // The other field is not used.
}
