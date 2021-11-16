// Package models contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID                     int
	Name                   string
	Password               string
	Email                  string
	FieldTypeA             SomeStructType // Mustn't be auto matched
	FieldSuperString       SuperString    // Can be auto matched with domain.Account.FieldSuperString cause SuperString is based on basic type `string`
	FieldReversedType      string
	FieldSliceReversedType []string
	SliceString            []string
}

// User represents the data model for a user.
type User struct {
	UserID   int
	Username string
	UserData UserData
}

// UserData represents data owned by the user.
// The fields of UserData operate at depth level 1.
type UserData struct {
	Options map[string]interface{}
	Data    Data
}

// Data represents a piece of data.
// The fields of UserData operate at depth level 2.
type Data struct {
	ID int
}

type SomeStructType struct {
	FieldB string
}

type SuperString string
