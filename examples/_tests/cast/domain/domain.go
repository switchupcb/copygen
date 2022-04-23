// Package domain contains business logic models.
package domain

type ReversedString string

// Account represents a user account.
type Account struct {
	ID             int
	Name           string
	Email          string
	SuperString    string
	ReversedString ReversedString
}
