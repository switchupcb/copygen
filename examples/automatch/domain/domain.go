// Package domain contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID                int
	Name              string
	User              DomainUser // The fields of DomainUser operate at depth level 1.
	Email             string
	FieldTypeA        SomeStructType // Mustn't be auto matched
	FieldSuperString  string         // Can be auto matched with models.Account.FieldSuperString cause SuperString is based on basic type `string`
	FieldReversedType ReversStringType
}

// DomainUser represents a user in relation to the business logic.
type DomainUser struct {
	UserID   int
	Username string
}

type SomeStructType struct {
	FieldA string
}

type ReversStringType string
