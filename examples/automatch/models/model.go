// The models package contains data storage models (i.e database).
package models

// Account represents the data model for account.
type Account struct {
	ID       int
	Name     string
	Password string
	Email    string
}

// A user represents the data model for a user.
type User struct {
	UserID   int
	Name     string
	UserData UserData // The fields of UserData operate at depth level 1.
}

type UserData struct {
	Options map[string]interface{}
	Data    Data // The fields of UserData operate at depth level 2.
}

type Data struct {
	ID int
}
