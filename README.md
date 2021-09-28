# Copygen

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge&logo=appveyor&logo=appveyor)](https://pkg.go.dev/github.com/switchupcb/copygen)
[![Go Report Card](https://goreportcard.com/badge/github.com/switchupcb/copygen?style=for-the-badge)](https://goreportcard.com/report/github.com/switchupcb/copygen)
[![MIT License](https://img.shields.io/github/license/switchupcb/copygen.svg?style=for-the-badge)](https://github.com/switchupcb/copygen/blob/main/LICENSE)

Copygen is a command-line [code generator](https://github.com/gophersgang/go-codegen) that generates type-to-type and field-to-field struct code without adding any reflection or dependencies to your project.

| Topic                           | Categories                                                                      |
| :------------------------------ | :------------------------------------------------------------------------------ |
| [Use](#use)                     | [Types](#types), [YML](#yml), [Command Line](#command-line), [Output](#output)  |
| [Customization](#customization) | [Templates](#templates), [Convert](#convert), [Custom Options](#options)        |
| [Matcher](#matcher)             | [Automatch](#automatch)                                                         |
| [Optimization](#optimization)   | [Shallow Copy vs. Deep Copy](#shallow-copy-vs-deep-copy), [Pointers](#pointers) |

### Benchmark

**The benefit to using Copygen is performance**: A benchmark by [gotidy/copy](https://github.com/gotidy/copy#benchmark) shows that a manual copy is **391x faster** than [jinzhu/copier](https://github.com/jinzhu/copier) and **3.97x faster** than the best reflection-based solution.

## Use

Each example has a **README**.

| Example                                                                         | Description                                                  |
| :------------------------------------------------------------------------------ | :----------------------------------------------------------- |
| main                                                                            | The default example _(README)_.                              |
| [automatch](https://github.com/switchupcb/copygen/tree/main/examples/automatch) | Uses the automatch feature _(doesn't require fields)_.       |
| [error](https://github.com/switchupcb/copygen/tree/main/examples/error)         | Uses templates to return an error (temporarily unsupported). |
| deepcopy                                                                        | Uses templates to create a deepcopy.                         |

This [example](https://github.com/switchupcb/copygen/blob/main/examples/main) uses three type-structs to generate the `ModelsToDomain()` function. All paths are specified from the `types.yml` filepath in `examples/main`.

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

A YML file is used to configure the code that is generated. View [Customization](#customization) for an example using custom templates. _Properties with `# default` are NOT necessary to include._

**types.yml**

```yml
# Define where the code will be generated.
generated:
  filepath: ./copygen.go
  package: copygen

# Define the imports that are included in the generated file.
import:
  - github.com/switchupcb/copygen/examples/main/domain
  - github.com/switchupcb/copygen/examples/main/models
  - github.com/switchupcb/copygen/examples/main/converter

# Define the functions to be generated.
# Properties with `# default` are NOT necessary to include.
functions:
 ModelsToDomain:

    # Custom function options can be defined for template use.
    options:                  # default: none
      custom: true

    # Define the types to be copied (to and from).
    to:
      Account:
        package:  domain      # default: none
        pointer:  true        # default: false     (# Optimization) 
  
        # Custom type options (to and from) can be defined for template use.
        options:              # default: none
          custom: false

    from:
      User:
        package: models       # default: none
        pointer: false        # default: false

        # Match fields to the to-type.
        fields:               # default: automatch (# Matcher)
          ID:
            to: UserID
            convert: c.Itoa   # default: none      (# Matcher)

            # Custom field options can be defined for template use.
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
// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.
package copygen

import (
	"github.com/switchupcb/copygen/examples/main/converter"
	"github.com/switchupcb/copygen/examples/main/domain"
	"github.com/switchupcb/copygen/examples/main/models"
)

// ModelsToDomain copies a User, Account to a Account.
func ModelsToDomain(tA *domain.Account, fU models.User, fA models.Account) {
	// Account fields
	tA.UserID = c.Itoa(fU.ID)
	tA.ID = fA.ID
	tA.Name = fA.Name

}
```

### Customization

The [error example](https://github.com/switchupcb/copygen/blob/main/examples/main) modifies the .yml to use **custom functions** which `return error`. This is done by modifying the .yml and creating custom template files. All paths are specified from the `types.yml` file path in `examples/error`.

**types.yml**

```yml
# Define where the code will be generated.
generated:
  filepath: ./copygen.go
  package: copygen

  # Define the optional custom templates used to generate the file.
  templates:
    header: ./templates/header.go
    function: ./templates/function.go
```

#### Templates

Templates can be created using **Go** to customize the generated code. The `copygen` generator uses the `package generator` `Header(*models.Generator)` to generate header code and `Function(*models.Function)` to generate code for each function. As a result, these _(package generator with functions)_ are **required** for your templates to work. View [models.Generator](https://github.com/switchupcb/copygen/blob/main/cli/models/function.go) and [models.Function](https://github.com/switchupcb/copygen/blob/main/cli/models/function.go) for context on the parameters passed to each function. Templates are interpreted by [yaegi](https://github.com/traefik/yaegi) which currently has limitations on module imports: As a result, **templates are temporarily unsupported.**

#### Convert

The `convert` property is used to specify a converter function. This is useful when you need to copy a value between two fields with different types, or provide another use-case. A default-template _converter function_ uses the following signature:

```go
func convert(f Field) Type {
    // implement custom function...
    return Type
}
```

where `Field` is replaced with the field it will receive _(i.e int)_, and `Type` is replaced with the type it will return _(i.e string)_.

#### Options

Function, Type, and Field custom options can be defined for template use. These are read into the respective [models](https://github.com/switchupcb/copygen/blob/main/cli/models) by [yaml](https://github.com/go-yaml/yaml).

## Matcher

Matching is specified in the `.yml` _(which functions as a schema in relation to other generators)_. All `from` type-fields are assigned to respective `to` types. The library assumes that it's used with other code generators: This complicates the use of tags which is why they aren't used.

### Automatch

If `fields` isn't specified for a `from` type, copygen will attempt to automatch type-fields by name. Automatch **supports field-depth** (where types are located within fields) **and recursive types** (where the same type is in another type). **You must specify the import path for types that use the automatcher.** Automatch loads types from Go modules _(in GOPATH)_. Ensure your modules are up to date by using `go get -u <insert/module/import/path>`.

#### Depth

A depth level of 0 will match the first-level fields. Increasing the depth level will match more fields.

```go
// depth level
type Account
  // 0
  ID      int
  Name    string
  Email   string
  Basic   domain.T // int
  User    domain.DomainUser
              // 1
              UserID   string
              Name     string
              UserData map[string]interface{}
  // 0
  Log     log.Logger
              // 1
              mu      sync.Mutex
                          // 2
                          state   int32
                          sema    uint32
              // 1
              prefix  string
              flag    int
              out     io.Writer
                          // 2
                          Write   func(p []byte) (n int, err error)
              buf     []byte
```

## Optimization 

### Shallow Copy vs. Deep Copy
The library generates a [shallow copy](https://en.m.wikipedia.org/wiki/Object_copying#Shallow_copy) by default. An easy way to deep-copy fields with the same return type is by using `new()` as/in a converter function or by using a custom template.

### Pointers
Go parameters are _pass-by-value_ which means that a parameter's value _(i.e int, memory address, etc)_ is copied into another location of memory. As a result, passing pointers to functions is more efficient **if the byte size of a pointer is less than the total byte size of the struct member's references**. However, be advised that doing so adds memory to the heap _[which can result in less performance](https://medium.com/@vCabbage/go-are-pointers-a-performance-optimization-a95840d3ef85)_. For more information regarding the use of pointers, read [Pointers vs. Values in Parameters and Return Values](https://stackoverflow.com/questions/23542989/pointers-vs-values-in-parameters-and-return-values/23551970#23551970). For more information on memory, read this article: [What Every Programmer Should Know About Memory](https://lwn.net/Articles/250967/).

## Contributing

You can contribute to this repository by viewing the [Project Structure, Code Specifications, and Roadmap](CONTRIBUTING.md).