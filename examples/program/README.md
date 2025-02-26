# Program

Using Copygen programmatically allows you to implement custom loaders, matchers, and generator functions, or pass `.yml` files similarly to the command line tool.

### Disclaimer

Copygen uses a [AGPLv3 License](../../README.md#license). 

An exception is provided for template and example files, which are licensed under the [MIT License](cli/generator/template/LICENSE.md). While this license lets you use the generator tool without restriction, **proprietary programmatic usage** of Copygen requires the purchase of a license exception. 

You can receive a license exception for Copygen by contacting SwitchUpCB using the [Copygen License Exception Inquiry Form](https://switchupcb.com/copygen-license-exception/).

## Standard

The `Run()` function operates similarly to using `copygen -yml path` in your command line interface.

### Define the Environment

Use an [Environment](https://pkg.go.dev/github.com/switchupcb/copygen/cli#Environment) object to specify the `.yml` Copygen will use.

```go
env := cli.Environment{
    YMLPath: "relative/path/to/yml/",
    Output:  false, // Don't output to standard output.
    Write:   true,  // Write the output to a file.
}
```

### Output

Generate code using `Run()`.

```go
code, err := env.Run()
if err != nil {
    return err
}

fmt.Println(code)
```

## Custom

The `Run()` function shows the flow of the standard `copygen` program.

```go
// The configuration file is loaded (.yml)
gen, err := config.LoadYML(e.YMLPath)
if err != nil {
    return "", fmt.Errorf("%w", err)
}

// The data file is parsed (.go)
if err = parser.Parse(gen); err != nil {
    return "", fmt.Errorf("%w", err)
}

// The matcher is run on the parsed data (to create the objects used during generation).
if err = matcher.Match(gen); err != nil {
    return "", fmt.Errorf("%w", err)
}

// The generator is used to generate code.
code, err := generator.Generate(gen, e.Output, e.Write)
if err != nil {
    return "", fmt.Errorf("%w", err)
}

return code, nil
```

Copygen's standard functions (`Parse`, `Match`, `Generate`) accept a [`model.Generator`](https://pkg.go.dev/github.com/switchupcb/copygen/cli/models#Generator). This lets the implementation of any configuration method (including `LoadYML`) with a Copygen standard function, as long as that method returns a valid `model.Generator`. similarly, custom `Parse`, `Match`, and `Generate` functions can be used with Copygen standard functions interchangeably.

### Configuration

The purpose of the config is to configure the `setup`, `output`, and `template` filepaths of a `Generator`; along with any `GeneratorOptions`. For more more information, read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/config#section-documentation).

### Parser

The purpose of the parser is to determine:

1. The code that is **kept** _(`Generator.Keep`)_ from the `setup` file _(placed above generated code in the `output` file)_.
2. Define the `Generator.Functions` (which contain fields).

For more information, read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/parser#section-documentation).

#### Parser Options

Read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/parser/options#section-documentation).

### Matcher

Read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/matcher#section-documentation).

### Generator

The generator package exports three methods of code generation by default. The easiest way to customize the generator is by using [templates](../../README.md#templates).

For more information, read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/generator#section-documentation).

## Debug

The [`debug.go`](https://pkg.go.dev/github.com/switchupcb/copygen/cli/models/debug#pkg-functions) file provides helper functions during debugging. Use a [GoAst Viewer](https://yuroyoro.github.io/goast-viewer/index.html) to view an Abstract Syntax Tree.

### PrintFieldGraph

**Parser**

```
type *models.Account
    Unpointed Field "*models.Account" of Definition "" Fields[4]: Parent ""
        Unpointed Field ID "*models.Account.ID" of Definition "int" Fields[0]: Parent "*models.Account"
        Unpointed Field Name "*models.Account.Name" of Definition "string" Fields[0]: Parent "*models.Account"
        Unpointed Field Password "*models.Account.Password" of Definition "string" Fields[0]: Parent "*models.Account"
        Unpointed Field Email "*models.Account.Email" of Definition "string" Fields[0]: Parent "*models.Account"
type *models.User
    Unpointed Field "*models.User" of Definition "" Fields[3]: Parent ""
        Unpointed Field UserID "*models.User.UserID" of Definition "int" Fields[0]: Parent "*models.User"
        Unpointed Field Name "*models.User.Name" of Definition "string" Fields[0]: Parent "*models.User"
        Unpointed Field UserData "*models.User.UserData" of Definition "string" Fields[0]: Parent "*models.User"
type *domain.Account
    Unpointed Field "*domain.Account" of Definition "" Fields[4]: Parent ""
        Unpointed Field ID "*domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account"
        Unpointed Field UserID "*domain.Account.UserID" of Definition "string" Fields[0]: Parent "*domain.Account"
        Unpointed Field Name "*domain.Account.Name" of Definition "string" Fields[0]: Parent "*domain.Account"
        Unpointed Field Other "*domain.Account.Other" of Definition "string" Fields[0]: Parent "*domain.Account"
```

**Matcher**

```
type *models.Account
    Unpointed Field "*models.Account" of Definition "" Fields[2]: Parent ""
        From Field ID "models.Account.ID" of Definition "int" Fields[0]: Parent "*models.Account"
        From Field Name "models.Account.Name" of Definition "string" Fields[0]: Parent "*models.Account"
type *models.User
    Unpointed Field "*models.User" of Definition "" Fields[1]: Parent ""
        From Field UserID "models.User.UserID" of Definition "int" Fields[0]: Parent "*models.User"
type *domain.Account
    Unpointed Field "*domain.Account" of Definition "" Fields[3]: Parent ""
        To Field ID "domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account"
        To Field UserID "domain.Account.UserID" of Definition "string" Fields[0]: Parent "*domain.Account"
        To Field Name "domain.Account.Name" of Definition "string" Fields[0]: Parent "*domain.Account"
```

_Use `PrintGeneratorFields` to call `PrintFunctionFields` on all of a generator's functions._

### PrintFieldTree

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

### PrintFieldRelation

```
To Field ID "*domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account" and From Field ID "models.Account.ID" of Definition "int" Fields[0]: Parent "*models.Account" are related to each other.
To Field ID "*domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account" is not related to From Field Name "models.Account.Name" of Definition "string" Fields[0]: Parent "*models.Account".
To Field ID "*domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account" is not related to From Field Password "models.Account.Password" of Definition "string" Fields[0]: Parent "*models.Account".
To Field ID "*domain.Account.ID" of Definition "int" Fields[0]: Parent "*domain.Account" is not related to From Field Email "models.Account.Email" of Definition "string" Fields[0]: Parent "*models.Account".
...
```