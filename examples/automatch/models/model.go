// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Password string
	Email    string
}

// User represents the data model for a user.
type User struct {
	UserID   int
	Name     string
	UserData UserData // The fields of UserData operate at depth level 1.
}

// UserData represents data owned by the user.
type UserData struct {
	Options map[string]interface{}
	Data    Data // The fields of UserData operate at depth level 2.
}

// Data represents a piece of data.
type Data struct {
	ID int
}
