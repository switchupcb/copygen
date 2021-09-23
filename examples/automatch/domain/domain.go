// The domain package contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID     int
	UserID string
	Name   string
	Email  string
}
