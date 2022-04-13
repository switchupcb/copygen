// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID       int    `api:"user_id"`
	Name     string `api:"name"`
	Email    string `other:"email"`
	Username string `api:"username" other:"tag"`
	Password string // The password field will not be copied.
}
