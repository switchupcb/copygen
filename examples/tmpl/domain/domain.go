// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID     int
	UserID string
	Name   string
	Other  string // The other field is not used.
}
