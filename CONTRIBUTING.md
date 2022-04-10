# Contributing

## Contributor License Agreement
 
Contributions to this project must be accompanied by a **Contributor License Agreement**. You or your employer retain the copyright to your contribution. This simply gives us permission to use and redistribute your contributions as part of the project.

## Pull Requests

Pull requests must pass all [CI/CD](#cicd) measures and follow the [code specification](#specification).

## Domain

The domain of Copygen lies in field manipulation. The program uses provided types to determine the fields we must assign. In this context, a "Type" refers to _the types used in a function (as parameters or results)_ rather than a type used to define variables. As the `parser` and `matcher` provides all required field information, you can improve Copygen by modifying the generator.

## Project Structure

The repository consists of a detailed [README](README.md), [examples](/examples/), and [**command line interface**](/cli/).

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

_Read [program](examples/program/README.md) for an overview of the application's code._


### Parser

A `setup` file's abstract syntax tree is traversed once, but involves four processes.

#### Keep

The `setup` file is parsed using an Abstract Syntax Tree. This tree contains the `type Copygen Interface` but also code that must be **kept** in the generated `output` file. For example, the package declaration, file imports, convert functions, and [custom types](README.md#custom-types) all exist _outside_ of the `type Copygen Interface`. Instead of storing these declarations and attempting to regenerate them, we simply discard declarations — from the `setup` file's AST — that won't be kept: In this case, the `type Copygen Interface` and `ast.Comments` (that refer to `Options`).

#### Options

**Convert** options are defined **outside** of the `type Copygen Interface` and may apply to multiple functions. As a result, all `ast.Comments` must be parsed before `models.Function` and `models.Field` objects can be created. In order to do this, the `type Copygen Interface` is stored, but **NOT** analyzed until the `setup` file is traversed. This leaves two ways to parse `ast.Comments` into `Options`.

1. Parse **Convert** `ast.Comments` into `Options` during `setup` file traversal, and **field** `ast.Comments` into `Options` _(defined above Copygen functions)_ while analyzing the `type Copygen Interface`.
2. Parse `ast.Comments` that were removed from the AST into `Options`.

Method **1** is slightly more efficient _(since **convert** `ast.Comments` are only referenced once; not stored)_, but only used because **convert** options require the name of their respective **convert** functions _(which can't be parsed from comments)_. In contrast, regex compilation is expensive — [especially in Go](https://github.com/mariomka/regex-benchmark#performance) — and avoided by only compiling unique comments once.

#### Imports

The `go/types` package will provide everything else; _**except**_ for alias import names. In order to assign aliased or non-aliased import names to `models.Field`, the imports of the `setup` file are mapped to a package path.

#### Copygen Interface

The `type Copygen interface` is parsed to setup the `models.Function` and `models.Field` objects used in the `Matcher` and `Generator`.
- [go/types Contents (Types, A -> B)](https://go.googlesource.com/example/+/HEAD/gotypes#contents)
- [go/packages Package Object](https://pkg.go.dev/golang.org/x/tools/go/packages#Package)
- [go/types Func (Signature)](https://pkg.go.dev/go/types#Func)
- [go/types Types](https://pkg.go.dev/go/types#pkg-types)

### Generator

Copygen supports three methods of generation for end-users _(developers)_: `.go`, `.tmpl`, and `programmatic`.

#### .go

`.go` code generation allows users to generate code using the programming language they are familiar with. `.go` code generation works by allowing the end-user to specify **where** _the `.go` file containing the code generation algorithm_ is, then running the file _at runtime_. In order to do this, we must use an **interpreter**. Templates are interpreted by our [temporary yaegi fork](https://github.com/switchupcb/yaegi). `models` objects are extracted via reflection and loaded into the interpreter. Then, the interpreter interprets the provided `.go` template file _(specified by the user)_ to run the `Generate()` function.

#### .tmpl

`.tmpl` code generation allows users to generate code using [`text/templates`](https://pkg.go.dev/text/template). `.tmpl` code generation works by allowing the end-user to specify **where** _the `.tmpl` file containing the code generation algorithm_ is, then parsing and executing the file _at runtime_.

#### programmatic

`programmatic` code generation allows users to generate code by using `copygen` as a third-party module. For more information, read the [program example README](/examples/program/README.md).

## Specification

### From vs. To

From and To is used to denote the direction of a type or field. A from-field is assigned **to** a to-field. In contrast, a from-field applies to all to-fields _(unless specified otherwise)_. As a result, **"From" comes before "To" when parsing** while **"To" comes before "From" in comparison**.

### Variable Names

| Variable | Description                                                                          |
| :------- | :----------------------------------------------------------------------------------- |
| from.*   | Variables preceded by from indicate from-functionality.                              |
| to.*     | Variables preceded by to indicate to-functionality.                                  |
| loadpath | `loadpath` represents the (relative) path of the loader (current working directory). |

### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

### Why Pointers

Contrary to the README, pointers aren't used — on Fields — as a performance optimization. Using pointers with Fields makes it less likely for a mistake to occur during the comparison of them. For example, using a for-copy loop on a `[]models.Field`:

```go
// A copy of field is created with a separate pointer.
for _, field := range fields {
   // fromField.To still points to the original field.
   // fromField.From points to a field which is NOT the copied field (but has the same values).
   if field == fromField.To {
      // Never Happens
      ...
   }
}
```

The same reasoning applies to `for i := 0; i < count; i++` loops.

### Anti-patterns

Using the `*models.Field` definition for a `models.Field`'s `Parent` field can be considered an anti-pattern. In the program, a `models.Type` specifically refers to the types in a function signature _(i.e `func(models.Account, models.User) *domain.Account`)_. While these types **are** fields _(which may contain other fields)_ , their actual `Type` properties are not relevant to `models.Field`. As a result, `models.Field` objects are pointed directly to maintain simplicity.

Using the `*models.Field` definition for a `models.Field`'s `From` and `To` fields can be placed into a `type FieldRelation`: `From` and `To` is only assigned in the matcher. While either method allows you to reference a `models.Field`'s respective `models.Field`, directly pointing `models.Field` objects adds more customizability to the program and more room for extension. 

## CI/CD

### Static Code Analysis

Copygen uses [golangci-lint](https://github.com/golangci/golangci-lint) in order to statically analyze code. You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.0` and run it using `golangci-lint run`. If you receive a `diff` error, you must add a `diff` tool in your PATH. There is one located in the `Git` bin.

If you receive `File is not ... with -...`, use `golangci-lint run --disable-all --no-config -Egofmt --fix`.

#### Fieldalignment

**Struct padding** aligns the fields of a struct to addresses in memory. The compiler does this to improve performance and prevent numerous issues on a system's architecture _(32-bit, 64-bit)_. As a result, misaligned fields add more memory-usage to a program, which can effect performance in a numerous amount of ways. For a simple explanation, view [Golang Struct Size and Memory Optimization](https://medium.com/techverito/golang-struct-size-and-memory-optimisation-b46b124f008d
). Fieldalignment can be fixed using the [fieldalignment tool](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment) which is installed using `go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest`.

**ALWAYS COMMIT BEFORE USING `fieldalignment -fix ./cli/...`** as it may remove comments.

### Tests

For information on testing, read [Integration Tests](examples/_tests/).

# Roadmap

Focus on these features:
   - Generate Templates: logic for all types in `.go` template + examples (in `tests`)
   - Generate Templates: stronger `.tmpl` template
   - Generator: deepcopy + example
   - CICD: workflow that ensures consistency between `generator/template` and `examples` template files.