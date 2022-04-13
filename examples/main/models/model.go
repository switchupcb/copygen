// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Password string
	Email    string
}

// A User represents the data model for a user.
type User struct {
	UserID   int
	Name     string
	UserData string
}
