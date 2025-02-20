# Contributing

## Contributor License Agreement
 
Contributions to this project must be accompanied by a **Contributor License Agreement**. 

You or your employer retain the copyright to your contribution: Accepting this agreement gives us permission to use and redistribute your contributions as part of the project.

## Pull Requests

Pull requests must pass all [CI/CD](#cicd) measures and follow the [code specification](#specification).

## Domain

The domain of Copygen lies in field manipulation. The program uses provided types to determine the fields we must assign. In the context of this domain, a `Type` refers to _the types used in a function (as parameters or results) [e.g., `func example(x TypeX) y TypeY`]_, and not a type used to define variables (e.g., `type example string`).

## Project Structure

The repository consists of a detailed [README](README.md), [examples](/examples/), and [**command line interface**](/cli/).

### Command Line Interface

The command-line interface _(cli)_ consists of 5 packages.

| Package   | Description                                                                                        |
| :-------- | :------------------------------------------------------------------------------------------------- |
| cli       | Contains the primary logic used to parse arguments and run the `copygen` command-line application. |
| models    | Contains models based on the application's functionality _(business logic)_.                       |
| config    | Contains external loaders used to configure the file settings and command line options.            |
| parser    | Uses an Abstract Syntax Tree (AST) and `go/types` to parse a setup file for fields.                |
| matcher   | Contains application logic to match fields to each other.                                          |
| generator | Contains the generator logic used to generate code _(and interpret templates)_.                    |

_Read [program](examples/program/README.md) for an overview of the application's code._


### Parser

A `setup` file's abstract syntax tree is traversed once, but involves four processes.

#### 1. Keep

The `setup` file is parsed using an Abstract Syntax Tree. 

This tree contains the `type Copygen Interface` but also code that must be **kept** in the generated `output` file. 

For example, the package declaration, file imports, convert functions, and [custom types](README.md#custom-types) all exist _outside_ of the `type Copygen Interface`. Instead of storing these declarations and attempting to regenerate them, we simply discard declarations — from the `setup` file's AST — that won't be kept: In this case, the `type Copygen Interface` and `ast.Comments` (that refer to `Options`).

#### 2. Options

**Convert** options are defined **outside** of the `type Copygen Interface` and may apply to multiple functions. So, all `ast.Comments` must be parsed before `models.Function` and `models.Field` objects can be created. In order to do this, the `type Copygen Interface` is stored, but **NOT** analyzed until the `setup` file is traversed. 

There are multiple ways to parse `ast.Comments` into `Options`, but **convert** options require the name of their respective **convert** functions _(which can't be parsed from comments)_. So, the most readable, efficient, and least error prone method of parsing `ast.Comments` into `Options` is to parse them when discovered and assign them from a `CommentOptionMap` later. In addition, regex compilation is expensive — [especially in Go](https://github.com/mariomka/regex-benchmark#performance) — and avoided by only compiling unique comments once.

#### 3. Copygen Interface

The `type Copygen interface` is parsed to setup the `models.Function` and `models.Field` objects used in the `Matcher` and `Generator`.
- [go/types Contents (Types, A -> B)](https://go.googlesource.com/example/+/HEAD/gotypes#contents)
- [go/packages Package Object](https://pkg.go.dev/golang.org/x/tools/go/packages#Package)
- [go/types Func (Signature)](https://pkg.go.dev/go/types#Func)
- [go/types Types](https://pkg.go.dev/go/types#pkg-types)

#### 4. Imports

The `go/types` package provides all of the other important information _**except**_ for alias import names. In order to assign aliased or non-aliased import names to `models.Field`, the imports of the `setup` file are mapped to a package path, then assigned to fields prior to matching.

### Generator

Copygen supports three methods of generation for end-users _(developers)_: `.go`, `.tmpl`, and `programmatic`.

#### .go

`.go` code generation allows users to generate code using the programming language they are familiar with. 

`.go` code generation works by allowing the end-user to specify **where** _the `.go` file containing the code generation algorithm_ is, then running the file _at runtime_. We must use an **interpreter** to provide this functionality.

`.go` templates are interpreted by a [yaegi fork](https://github.com/switchupcb/yaegi). 
1. `models` objects are extracted via reflection and loaded into the interpreter. 
2. Then, the interpreter interprets the provided `.go` template file _(specified by the user)_ to run the `Generate()` function.

#### .tmpl

`.tmpl` code generation allows users to generate code using [`text/templates`](https://pkg.go.dev/text/template). 

`.tmpl` code generation works by allowing the end-user to specify **where** _the `.tmpl` file containing the code generation algorithm_ is, then parsing and executing the file _at runtime_.

#### programmatic

`programmatic` code generation lets users generate code by using `copygen` as a third-party module. For more information, read the [program example](/examples/program/README.md).

## Specification

### From vs. To

From and To is used to denote the direction of a type or field. A from-field is assigned **to** a to-field. In contrast, one from-field can match many to-fields. So, **"From" comes before "To" when parsing** while **"To" comes before "From" when matching**.

### Variable Names

| Variable | Description                                             |
| :------- | :------------------------------------------------------ |
| from.*   | Variables preceded by from indicate from-functionality. |
| to.*     | Variables preceded by to indicate to-functionality.     |

### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

### Why Pointers

Contrary to the README, pointers aren't used — on `models.Fields` — as a performance optimization. Using pointers with `models.Fields` makes it less likely for a mistake to occur during their comparison. For example, using a for-copy loop on a `[]models.Field`:

```go
// A copy of field is created with a distinct memory address.
for _, field := range fields {
   // field.To still points to the original field's .To memory address.
   // field.To.From points to the original field's memory address, which is NOT the copied field's memory address, even though both fields' fields have the same values.
   if field == field.To.From {
      // never happens
      ...
   }
}
```

### Anti-patterns

Using the `*models.Field` definition for a `models.Field`'s `Parent` field can be considered an anti-pattern. In the program, a `models.Type` specifically refers to the types in a function signature _(i.e `func(models.Account, models.User) *domain.Account`)_. While these types **are** fields _(which may contain other fields)_ , their actual `Type` properties are not relevant to `models.Field`. So, `models.Field` objects are pointed directly to each other for simplicity.

Using the `*models.Field` definition for a `models.Field`'s `From` and `To` fields can be placed into a `type FieldRelation` since `From` and `To` is only assigned in the matcher. While either method allows you to reference a `models.Field`'s respective `models.Field`, directly pointing `models.Field` objects adds more customizability to the program for the end user.

## CI/CD

### Static Code Analysis

Copygen uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.1` and run it using `golangci-lint run`. If you receive a `diff` error, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

#### Fieldalignment

**Struct padding** aligns the fields of a struct to addresses in memory. The compiler does this to improve performance and prevent numerous issues on a system's architecture _(32-bit, 64-bit)_. So, misaligned fields add more memory-usage to a program, which can effect performance in a numerous amount of ways. For a simple explanation, view [Golang Struct Size and Memory Optimization](https://medium.com/techverito/golang-struct-size-and-memory-optimisation-b46b124f008d
).

Fieldalignment can be fixed using the [fieldalignment tool](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) which is installed using `go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest`.

**ALWAYS COMMIT BEFORE USING `fieldalignment -fix ./cli/...`** as it may remove comments.

### Tests

For information on testing, read [Tests](examples/_tests/).

# Roadmap

Implement the following features.
   - Generator: deepcopy
   - Parser: Fix Free-floating comments _(add structs in [`multi`](examples/_tests/multi/copygen.go) to test)_
