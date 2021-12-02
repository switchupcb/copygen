// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID               int
	Name             string
	User             DomainUser // The fields of DomainUser operate at depth level 1.
	Email            string
	Balance          float64 `api:"Bal"`
	FieldSuperString string  // Can be auto matched with models.Account.FieldSuperString cause SuperString is based on basic type `string`
	SliceString      []string
	TagName          string `modelname:"StrangeName"`
}

// DomainUser represents a user in relation to the business logic.
type DomainUser struct {
	UserID   int
	Username string
}
