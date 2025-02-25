# Example: Cast

These examples use the `cast from to modifier` option to modify fields before assignment: This lets you perform direct type assertion, conversion, expressions, function usage, and property usage with a matched field.

| Example                                 | Description                                                            |
| :-------------------------------------- | :--------------------------------------------------------------------- |
| [Assert](examples/cast/assert/)         | Use `cast` generator options to enable automatic type assertion.       |
| [Convert](examples/convert/)            | Use `cast` function option modifiers to enable direct type conversion. |
| [Depth](examples/cast/depth/)           | Use `cast` option modifier `depth` to change autocasting behavior.     |
| [Expression](examples/cast/expression/) | Use `cast` option modifier `expression` to evaluate an expression.     |
| [Function](examples/cast/function/)     | Use `cast` option modifier `func()` to call a function.                |
| [Property](examples/cast/property/)     | Use `cast` option modifier `.Property` to reference a type property.   |

## Modifiers

Copygen handles assignability of fields with identical field names and definitions by default. Enable Copygen type casting when you must match fields using identical field names, but different definitions: This option modifies the matcher to match fields using an automatic or provided modifier.

_Typecasting is not a feature in the Go programming language._

### Assertion

[Type assertion](https://go.dev/ref/spec#Type_assertions) provides access to an interface's underlying concrete values: This lets you use the underlying concrete type's fields and functions. Enabling **type assertion** in Copygen applies to the following cases.
- Assignment of objects to interfaces: `interface = type`.
- Assertion of interfaces into objects: `type = interface.(type)`.

### Conversion

[Type conversion](https://go.dev/ref/spec#Conversions) changes the type of an expression to the type specified by the conversion: This lets you convert types into other types. Enabling **type conversion** in Copygen enables type conversion at the specified **depth**.

_For example, `custom = custom(bool)` represents a type conversion at a depth of one._

### Expressions

An expression is a value: Multiple expressions can be used to create statements which perform operations. Copygen `cast` modifiers are copied directly from the option. So using literal expressions _(i.e `* 2`)_, type functions _(i.e `.String()`)_, and type properties _(i.e `.Property`)_ is allowed.

_For example, `cast model.Type.Field domain.Type.Field * 2` results in the assignment `model.Type.Field = domain.Type.Field * 2` upon a match of the fields._

## Usage

There are multiple ways to enable Copygen type casting.

### Generator Options

```yml
# Define how the matcher will work.
matcher:
  # Skip the matcher (default: false).
  skip: false

  # Control the matcher type cast behavior.
  cast:
    # Enable automatic casting (default: false).
    enabled: true
    
    # Set the maximum depth for automatic casting (default: 1)
    depth: 1  
    
    # Disable certain features of casting.
    disabled:

      # Disable assignment of objects to interfaces (default: false).
      assignObjectInterface: false

      # Disable assertion of interfaces to objects (default: false).
      assertInterfaceObject: false

      # Disable type conversion (default: false).
      convert: false
```

### Function Option

Use the `cast from to modifier` function option to enable casting for the respective function. Regex is supported for from-fields. The **modifier** flag is optional.
- `cast .* package.Type.Field`
- `cast models.Account.ID domain.Account.ID .String()`

### Option Modifier

Use the `-cast modifier` function option modifier to enable casting for the respective function option. The **modifier** flag is optional.
- `automatch package.Type.Field -cast`
- `automatch models.User.* -cast 1`
- `map .* package.Type.Field -cast .String()`
- `map .* package.Type.Field -cast Convert()`
- `tag package.Type.Field key -cast .Property`
- `tag .* api -cast + 100`

## Behavior

The `cast` option is a **modifier**: It modifies the matching algorithm or assignment of fields. It can't be used to match fields in a direct manner _(unlike `automatch`, `map`, and `tag`)_.

Copygen will perform **automatic typecasting** using type assertion or conversion at the specified depth level when a `modifier` is **NOT** provided. Otherwise, the definition of the **provided modifier** is evaluated to match fields. 

_For example, `map .* package.Type.Field -cast .String()` matches from-fields with the name `Field` (from `package.Type.Field`) and definition `string` (since `.String()` returns `string`) when depth is greater than 0._