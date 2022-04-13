// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// User represents the data model for a user.
type User struct {
	UserID   int
	Username string
	UserData UserData // The fields of UserData operate at depth level 2.
}

// UserData represents data owned by the user.
type UserData struct {
	Options map[string]interface{}
	Data    Data // The fields of Data operate at depth level 3.
}

// Data represents a piece of data.
type Data struct {
	ID int
}
