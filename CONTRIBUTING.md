# Contributing

## Contributor License Agreement
 
Contributions to this project must be accompanied by a **Contributor License Agreement**. You or your employer retain the copyright to your contribution. This simply gives us permission to use and redistribute your contributions as part of the project.

## Pull Requests

Pull requests must work with all examples _(confirmed through [integration tests]())_ and follow the [code specification](#guideline).

## Domain

The domain of Copygen lies in field manipulation. The program uses provided types to determine the fields we must assign. In this context, a "Type" refers to _the types used in a function (as parameters or results)_ rather than a type used to define variables. As the `parser` and `matcher` provides all required field information, you can improve Copygen by modifying the generator.

## Project Structure

The repository consists of a detailed [README](https://github.com/switchupcb/copygen#copygen), [examples](https://github.com/switchupcb/copygen/tree/main/example), and [**command line interface**](https://github.com/switchupcb/copygen).

### Command Line Interface

The command-line interface _(cli)_ consists of 5 packages. 

| Package   | Description                                                                                        |
| :-------- | :------------------------------------------------------------------------------------------------- |
| cli       | Contains the primary logic used to parse arguments and run the `copygen` command-line application. |
| models    | Contains models based on the application's functionality _(business logic)_.                       |
| config    | Contains external loaders used to configure the file settings and command line options.            |
| parser    | Uses Abstract Syntax Tree (AST) analysis to parse a data file for fields.                          |
| matcher   | Contains application logic to match fields to each other _(manually and automatically)_.           |
| generator | Contains the generator logic used to generate code _(and interpret templates)_.                    |

The command line interface package allows you see the flow of the program.
```go
// The configuration file is loaded (.yml)
gen, err := config.LoadYML(e.YMLPath)
if err != nil {
    return err
}

// The data file is parsed (.go)
if err = parser.Parse(gen); err != nil {
    return err
}

// The matcher is run on the parsed data (to create the objects used during generation).
if err = matcher.Match(gen); err != nil {
    return err
}

// The generator is used to generate code.
if err = generator.Generate(gen, e.Output); err != nil {
    return err
}

return nil
```

#### Parser

A setup file's Abstract Syntax Tree is traversed once. This is done in three steps:
1. **Options:** Regex compilation is expensive — [especially in Go](https://github.com/mariomka/regex-benchmark#performance) — and avoided by only compiling unique option-comments once. The location of a `convert` option cannot be assumed: Therefore, we must traverse the entire Abstract Syntax Tree in order to correctly assign options. As a result, the `type Copygen Interface` is stored for post-traversal analysis.
2. **Keep:** The code that is kept after generation is stored — or more so kept — in the AST. We do not want to keep option-comments nor the Copygen interface in the AST. However, they must still be present during the `type Copygen Interface` analysis _(which requires the option-comments)_. As a result, comments are stored in the parser for post-analysis removal.
3. **type Copygen Interface:** The `type Copygen interface` is parsed to setup the function and fields used in the program. 

### Specification

#### From vs. To

From and To is used to denote the direction of a type or field. A from-field is assigned **to** a to-field. In contrast, a from-field applies to all to-fields _(unless specified otherwise)_. As a result, **"From" comes before "To" when parsing** while **"To" comes before "From" in comparison**.

#### Variable Names

| Variable | Description                                                                          |
| :------- | :----------------------------------------------------------------------------------- |
| from.*   | Variables preceded by from indicate to-functionality.                                |
| to.*     | Variables preceded by to indicate to-functionality.                                  |
| loadpath | `loadpath` represents the (relative) path of the loader (current working directory). |

#### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

#### Why Pointers

Contrary to the README, pointers aren't used — on Fields — as a performance optimization. Using pointers with Fields makes it less likely for a mistake to occur during the comparison of them. For example, using a for-copy loop on a `[]models.Field`:

```go
// A copy of field is created with a separate pointer.
for _, field := range fields {
   // field.To still points to the original field.
   // fromField.From points to a field which is NOT the copied field (but has the same values)
   if field == fromField.To {
      // Never Happens
      ...
   }
}
```

The same reasoning applies to `for i := 0; i < count; i++` loops.

#### Anti-patterns

Using the `*models.Field` definition for a `models.Field`'s `Parent` field can be considered an anti-pattern. In the program, a `models.Type` specifically refers to the types in a function signature _(i.e `func(models.Account, models.User) *domain.Account`)_. While these types **are** fields _(which may contain other fields)_ , their actual `Type` properties are not relevant to `models.Field`. As a result, `models.Field` objects are pointed directly to maintain simplicity.

Using the `*models.Field` definition for a `models.Field`'s `From` and `To` fields can be placed into a `type FieldRelation`: `From` and `To` is only assigned in the matcher. While either method allows you to reference a `models.Field`'s respective `models.Field`, directly pointing `models.Field` objects adds more customizability to the program and more room for extension.

#### CI/CD

Copygen uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1` and run it using `golangci-lint run`. If you receive a `diff` error, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

# Roadmap

Focus on these features:
   - Generator: deepcopy + example
   - Integration Tests (`examples`)

## Debug

The `debug.go` provides helper functions during debugging. To view an Abstract Syntax Tree, use a [GoAst Viewer](https://yuroyoro.github.io/goast-viewer/index.html).

#### PrintFunctionFields

**Parser**

```
type models.Account.
      Unpointed Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0]
      Unpointed Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0]
      Unpointed Field "models.Account.Password" of Definition "string": Parent "models.Account." Fields[0]
      Unpointed Field "models.Account.Email" of Definition "string": Parent "models.Account." Fields[0]
type models.User.
        Unpointed Field "models.User.ID" of Definition "int": Parent "models.User." Fields[0]
        Unpointed Field "models.User.Name" of Definition "int": Parent "models.User." Fields[0]
        Unpointed Field "models.User.UserData" of Definition "string": Parent "models.User." Fields[0]
type *domain.Account.
        Unpointed Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0]
        Unpointed Field "*domain.Account.UserID" of Definition "string": Parent "*domain.Account." Fields[0]
        Unpointed Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0]
        Unpointed Field "*domain.Account.Other" of Definition "string": Parent "*domain.Account." Fields[0]
```

**Matcher**

```
type models.Account.
        From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0]
        From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0]
type models.User.
        From Field "models.User.UserID" of Definition "int": Parent "models.User." Fields[0]
type *domain.Account.
        To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0]
        To Field "*domain.Account.UserID" of Definition "int": Parent "*domain.Account." Fields[0]
        To Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0]
```

#### PrintFieldGraph

**Parser**

```
Unpointed Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0]
Unpointed Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0]
Unpointed Field "models.Account.Password" of Definition "string": Parent "models.Account." Fields[0]
Unpointed Field "models.Account.Email" of Definition "string": Parent "models.Account." Fields[0]
Unpointed Field "models.User.ID" of Definition "int": Parent "models.User." Fields[0]
Unpointed Field "models.User.Name" of Definition "int": Parent "models.User." Fields[0]
Unpointed Field "models.User.UserData" of Definition "string": Parent "models.User." Fields[0]
Unpointed Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0]
Unpointed Field "*domain.Account.UserID" of Definition "string": Parent "*domain.Account." Fields[0]
Unpointed Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0]
Unpointed Field "*domain.Account.Other" of Definition "string": Parent "*domain.Account." Fields[0]
```

**Matcher**

```
From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0]
From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0]
From Field "models.User.UserID" of Definition "int": Parent "models.User." Fields[0]
To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0]
To Field "*domain.Account.UserID" of Definition "string": Parent "*domain.Account." Fields[0]
To Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0]
```

#### PrintFieldTree

**Parser** 

```go
type Account // domain
    // 0
    ID      int
    Name    string
    Email   string
            // 1
    User    domain.DomainUser
                UserID  int
                Name    string    
type User    // models
    // 0 
    UserID    int
    Name      string
              // 1
    UserData  models.UserData
                  Options map[string]interface{}
                  // 2
                  Data    models.Data
                        ID      int
type Account // models
    // 0
    ID       int
    Name     string
    Password string
    Email    string
```

**Matcher**

```go
// depth-level 0 tree
type Account
    ID      int
    Name    string
type User
    UserID  int
type Account
    ID      int
    UserID  int
    Name    string
```

#### PrintFieldRelation

**Matcher (Unpointed)**

```
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] and From Field From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0] are related to each other.
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field Unpointed Field "models.Account.Password" of Definition "string": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field Unpointed Field "models.Account.Email" of Definition "string": Parent "models.Account." Fields[0].
To Field Unpointed Field "*domain.Account.UserID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0].
...
```

**Matcher (Pointed)**

```
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] and From Field From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0] are related to each other.
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.UserID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.UserID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.Account.ID" of Definition "int": Parent "models.Account." Fields[0].
To Field To Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0] and From Field From Field "models.Account.Name" of Definition "string": Parent "models.Account." Fields[0] are related to each other.
To Field To Field "*domain.Account.ID" of Definition "int": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.User.UserID" of Definition "int": Parent "models.User." Fields[0].
To Field To Field "*domain.Account.UserID" of Definition "int": Parent "*domain.Account." Fields[0] and From Field From Field "models.User.UserID" of Definition "int": Parent "models.User." Fields[0] are related to each other.
To Field To Field "*domain.Account.Name" of Definition "string": Parent "*domain.Account." Fields[0] is not related to From Field From Field "models.User.UserID" of Definition "int": Parent "models.User." Fields[0].
```

#### CountFields

```
6
```