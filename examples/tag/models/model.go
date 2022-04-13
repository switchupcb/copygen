// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int    `api:"id"`
	Name     string `api:"name"`
	Password string
	Email    string `other:"email"`
}

// User represents the data model for a user.
type User struct {
	UserID   int    `api:"user_id"`
	Username string `api:"username,omitempty"`
}
