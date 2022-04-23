// Package models contains data storage models (i.e database).
package models

type SuperString string

// Account represents the data model for account.
type Account struct {
	ID             int
	Name           string
	Password       string
	Email          string
	SuperString    SuperString
	ReversedString string
}
