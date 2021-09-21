# Copygen

Copygen is a [Go code generator](https://github.com/gophersgang/go-codegen) that uses [Jennifer](https://github.com/dave/jennifer) in order to generate type-to-type and field-to-field struct code without reflection.

## Benchmark

**The benefit to using Copygen is performance**: A benchmark by [gotidy/copy](https://github.com/gotidy/copy#benchmark) shows that a manual copy is **391x faster** than [jinzhu/copier](https://github.com/jinzhu/copier) and **3.97x faster** than the best reflection-based solution.

![copy-benchmark](https://image.prntscr.com/image/-AcdCKSQSiqmrJ4KAW_ODg.png)

## Use

This example uses two type-structs (`INSERT NAME` and `INSERT NAME 2`) to generate the `ModelsToDomain()` function.

### Structs

`path/to/domain`

```go
package domain
// INSERT STRUCT CODE
```

`path/to/models`
```go
package models
// INSERT STRUCT CODE
```

### YML

A YML file is used to configure the code that is generated.

**struct.yml**

```yml
# Define where the code will be generated.
generated:
  name: INSERT EXAMPLE COPYGEN PATH
  package: copygen

# Define the imports that are included in the generated file.
import:
  - INSERT EXAMPLE DOMAIN PACKAGE
  - INSERT EXAMPLE MODEL PACKAGE
  - INSERT EXAMPLE CONVERTER PACKAGE c

# Define the function(s) to be generated.
# Properties with `# default` are NOT necessary to include.
function:
  name: ModelsToDomain

  # Define the types to be copied (to and from).
  # Note: Type-properties (i.e 'struct') can have any name.
  types:
    to:
      struct: INSERT DOMAIN
        filename: INSERT PATH TO DOMAIN
        pointer:  true        # default: false  (Optimization)
        deepcopy: false       # default: false  (Optimization)
       
    from:
      struct: INSERT MODEL
        filename: INSERT PATH TO MODEL

        # Match fields to the to-type.
        fields:
          INSERT EXAMPLE FIELD:
            to: INSERT RESPECTIVE TO-FIELD
            convert: c.itoa   # default: none  (Matcher)
          INSERT EXAMPLE FIELD2:
            to: INSERT RESPECTIVE TO-FIELD.
            convert: new

      struct: INSERT MODEL2
        filename: INSERT PATH TO MODEL2
        fields:
          INSERT EXAMPLE FIELD:
            to: INSERT RESPECTIVE TO-FIELD
```

_See [Optimization]() or [Matcher]() for information on respective properties._

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

### Output

This will output a `copygen.go` file _(defined in the .yml above)_ with the specified functions.

```go
// INSERT OUTPUT RESULT
```

 View the [example]() or [tests]().

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