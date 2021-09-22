# Copygen

Copygen is a [Go code generator](https://github.com/gophersgang/go-codegen) that generates type-to-type and field-to-field struct code without reflection.

**Topics**

| Topic                         | Categories                                                                      |
| :---------------------------- | :------------------------------------------------------------------------------ |
| [Use](#use)                   | [Types](#types), [YML](#yml), [Command Line](#command-line), [Output](#output)  |
| [Matcher](#matcher)           | [Convert](#convert)                                                             |
| [Optimization](#optimization) | [Shallow Copy vs. Deep Copy](#shallow-copy-vs-deep-copy), [Pointers](#pointers) |

## Benchmark

**The benefit to using Copygen is performance**: A benchmark by [gotidy/copy](https://github.com/gotidy/copy#benchmark) shows that a manual copy is **391x faster** than [jinzhu/copier](https://github.com/jinzhu/copier) and **3.97x faster** than the best reflection-based solution.

![copy-benchmark](https://image.prntscr.com/image/-AcdCKSQSiqmrJ4KAW_ODg.png)

## Use

This [example](https://github.com/switchupcb/copygen/blob/main/example) uses three type-structs to generate the `ModelsToDomain()` function. All paths are specified from the `types.yml` file path in `examples`.

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

### YML

A YML file is used to configure the code that is generated.

**types.yml**

```yml
# Define where the code will be generated.
generated:
  filepath: ./copygen.go
  package: copygen

# Define the imports that are included in the generated file.
import:
  - github.com/switchupcb/copygen/domain
  - github.com/switchupcb/copygen/models
  - github.com/switchupcb/copygen/converter

# Define the functions to be generated.
# Properties with `# default` are NOT necessary to include.
functions:

 # Custom function options can be defined for template use.
 ModelsToDomain:
    options:                  # default: none
      custom: true

    # Define the types to be copied (to and from).
    # Custom type options (to and from) can be defined for template use.
    to:
      Account:
        package:  domain      # default: none
        pointer:  true        # default: false  (Optimization) 
        options:              # default: none
          custom: false

    from:
      User:
        package: models       # default: none
        pointer: false        # default: false

        # Match fields to the to-type.
        # Custom field options can be defined for template use.
        fields:
          ID:
            to: UserID
            convert: c.Itoa   # default: none  (Matcher)
            options:          # default: none
              custom: false

      Account:
        package: models       # default: none
        fields:
          ID:
            to: ID
          Name:
            to: Name
```

_See [Optimization](https://github.com/switchupcb/copygen#Optimization) or [Matcher](https://github.com/switchupcb/copygen#matcher) for information on respective properties._

### Command Line

Install the command line utility _(outside your project)_.

```
go install github.com/switchupcb/copygen
```

Run with given options.

```bash
# Specify the yml file.
copygen -yml path/to/yml
```

_The path to the YML file is specified in reference to the current working directory._

### Output

This example outputs a `copygen.go` file with the specified imports and functions.

```go
package copygen

import (
	"github.com/switchupcb/copygen/example/converter"
	"github.com/switchupcb/copygen/example/domain"
	"github.com/switchupcb/copygen/example/models"
)

// ModelsToDomain copies fields from a models User and models Account to a domain Account.
func ModelsToDomain(tA *domain.Account, fU models.User, fA models.Account) error {
	tA.UserID = c.Itoa(fU.ID)
	tA.ID = fA.ID
	tA.Name = fA.Name
	return nil
}
```

## Matcher

Matching is specified in the `.yml` _(which functions as a schema in relation to other generators)_. This library assumes that it's used with other code generators which would make using tags difficult. We also avoid automatic matching _(by name or position)_ because there are few cases where it's viable.

### Convert

The `convert` property is used to specify a converter function. This is useful when you need to copy a value between two fields with different types, or provide another use-case. A _converter function_ uses the following signature:

```go
func convert(f Field) Type {
    // implement custom function...
    return Type
}
```

where `Field` is replaced with the field it will receive _(i.e int)_, and `Type` is replaced with type it will return _(i.e string)_.

## Optimization 

### Shallow Copy vs. Deep Copy
The library generates a [shallow copy](https://en.m.wikipedia.org/wiki/Object_copying#Shallow_copy) by default. An easy way to deep-copy fields with the same return type is by using `new()` as/in a converter function.

### Pointers
Go parameters are _pass-by-value_ which means that a parameter's value _(i.e int, memory address, etc)_ is copied into another location of memory.

As a result, passing pointers to functions is more efficient **if the byte size of a pointer is less than the total byte size of the struct member's references**. However, be advised that doing so adds memory to the heap _[which can result in less performance](https://medium.com/@vCabbage/go-are-pointers-a-performance-optimization-a95840d3ef85)_. 

You can read this article for more information on memory: [What Every Programmer Should Know About Memory](https://lwn.net/Articles/250967/).

## Contributing

You can contribute to this repository by viewing the [Project Structure](CONTRIBUTING.md).
