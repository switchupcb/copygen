# Example: Manual

The manual examples uses manual mapping to match three models.


### Types

`./domain/domain.go`

```go
// The domain package contains business logic models.
package domain

// Account represents a user account.
type Account struct {
	ID     int
	UserID string
	Name   string
	Other  string // The other field is not used.
}
```

`./models/model.go`

```go
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
	ID       int
	Name     int
	UserData string
}
```

### Setup

Setting up copygen is a 2-step process involving a `YML` and `GO` file.

**setup.yml**

```yml
# Define where the code will be generated.
generated:
  setup: ./setup.go
  output: ./copygen.go
  package: copygen

# Templates and custom options aren't used for this example.
```

**setup.go**

Create an interface in the specified setup file with a `type Copygen interface`. In each function, specify _the types you want to copy from_ as parameters, and _the type you want to copy to_ as return values.

```go
/* Copygen defines the functions that will be generated. */
type Copygen interface {
	// map models.Account.ID domain.Account.ID
	// map models.Acount.Name domain.Account.Name
	// map models.User.ID domain.Account.UserID
	// alloc
	ModelsToDomain(models.Account, models.User) *domain.Account
}

/* Define the fields this converter is applied to using regex. CONVERTERS ARE ONLY APPLIED TO VALID FIELDS. */
// convert: models.User.ID
// comment: Itoa converts an integer to an ascii value.
func Itoa(i int) string {
	return strconv.Itoa(i)
}
```

## Output

`copygen -yml path/to/yml`

```go
INSERT
```