# Program

Using Copygen programmatically allows you to implement custom loaders, matchers, and generator functions, or pass `.yml` files in a similar manner to the command line tool.

### Disclaimer

Copygen uses a [GPLv3 License](../../README.md#license). An exception is provided for template and example files, which are licensed under the [MIT License](cli/generator/template/LICENSE.md). While this allows you to use the generator tool without restriction, **proprietary programmatic usage** of Copygen requires the purchase of a license exception. In order to purchase a license exception, please contact SwitchUpCB using the [Copygen License Exception Inquiry Form](https://switchupcb.com/copygen-license-exception/). For more information, please read [What Can I do?](../../README.md#what-can-i-do)

## Standard

The `Run()` function operates in a similar manner to using `copygen -yml path` in your command line interface.

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

Copygen's standard functions (`Parse`, `Match`, `Generate`) accept a [`model.Generator`](https://pkg.go.dev/github.com/switchupcb/copygen/cli/models#Generator). This allows the implementation of any configuration method (including `LoadYML`) with a Copygen standard function, as long as that method returns a valid `model.Generator`. In a similar manner, custom `Parse`, `Match`, and `Generate` functions can be used with Copygen standard functions interchangeably.

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

The generator package exports three methods of code generation by default. The easiest way to customize the generator is by using [templates](../../README.md#templates). For more information, read the [documentation](https://pkg.go.dev/github.com/switchupcb/copygen/cli/generator#section-documentation).


## Debug

The [`debug.go`](https://pkg.go.dev/github.com/switchupcb/copygen/cli/models/debug#pkg-functions) file provides helper functions during debugging. To view an Abstract Syntax Tree, use a [GoAst Viewer](https://yuroyoro.github.io/goast-viewer/index.html).

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