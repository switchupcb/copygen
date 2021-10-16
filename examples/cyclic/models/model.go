// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Password string
	Info     User // User matches with the `Info Cyclic` field.
}

// A User represents the data model for a user.
type User struct {
	UserID   int
	Username string
}
