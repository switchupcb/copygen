# Contributing

## License

You agree to license any contribution to this library under the [MIT License](https://github.com/switchupcb/copygen/blob/main/LICENSE).

## Pull Requests

Pull requests must work with all examples _(confirmed through [integration tests]())_ and follow the [code specification](#guideline).

## Domain

 The domain of copygen lies in field manipulation. The program uses provided types to determine the fields we must assign. In this context, a "Type" refers to _the types used in a function (as parameters)_ rather than a type used to define variables. As the `matcher/ast` provides all required field information, you can improve copygen by modifying the generator.

 ### Improving the Generator

The generator interprets the provided template file at runtime to generate code. The generator can be improved by adding custom deepcopy functions for slice, map, array, etc that the user can use (via options).

## Project Structure

The repository consists of a detailed [README](https://github.com/switchupcb/copygen#copygen), [examples](https://github.com/switchupcb/copygen/tree/main/example), and [**command line interface**](https://github.com/switchupcb/copygen).

### Command Line Interface

The command-line interface _(cli)_ consists of 4 packages. 

| Package   | Description                                                                                      |
| :-------- | :----------------------------------------------------------------------------------------------- |
| cli       | Contains the primary logic used to parse arguments and run the copygen command-line application. |
| models    | Contains models based on the application's functionality _(logic)_.                              |
| loader    | Contains external loaders used to configure the code that is generated.                          |
| generator | Contains the generator logic used to generate code _(and interpret templates)_.                  |

### Specification

#### To vs. From

To and From is used to denote the direction of a type or field. A to-field is assigned **from** a from-field. Vice-versa, _all_ from-fields are assigned **to** to-fields: Therefore, **"To" always comes before "From"**.

#### Variable Names

**Common Variable Names**

| Variable | Description                                                                          |
| :------- | :----------------------------------------------------------------------------------- |
| to.*     | Variables preceded by to indicate to-functionality.                                  |
| from.*   | Variables preceded by from indicate to-functionality.                                |
| loadpath | `loadpath` represents the (relative) path of the loader (current working directory). |

#### Comments

Comments follow [Effective Go](https://golang.org/doc/effective_go#commentary) and explain why more than what _(unless the "what" isn't intuitive)_.

#### Why Pointers

Contrary to the README, pointers aren't used — on Fields in the matcher — as a performance optimization. Using pointers with Fields makes it less likely for a mistake to occur during the comparison of them. For example, using a for-copy loop on a `[]models.Field`:

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

# Roadmap

Focus on these features:
   - Generator: deepcopy + example
   - Integration Tests (`examples`)
   - Workflows for Pull Requests (Static Code Analysis)

## Debug

The `debug.go` provides helper functions during debugging.

#### PrintFunctionFields

```
TO type domain.Account
   To Field Name of Definition string: Parent 0xc000234200 Field OF 0x0 Fields []        
   To Field UserID of Definition int: Parent 0xc000234400 Field OF 0xc0002be500 Fields []
   To Field ID of Definition int: Parent 0xc000234500 Field OF 0x0 Fields []
   To Field Name of Definition string: Parent 0xc000234600 Field OF 0x0 Fields []        
   To Field Email of Definition string: Parent 0xc000234700 Field OF 0x0 Fields []
FROM type models.User
   From Field Name of Definition string: Parent 0xc000234800 Field OF 0x0 Fields []
   From Field UserID of Definition int: Parent 0xc000234900 Field OF 0x0 Fields []
FROM type models.Account
   From Field ID of Definition int: Parent 0xc000234a00 Field OF 0x0 Fields []
   From Field Name of Definition string: Parent 0xc000234b00 Field OF 0x0 Fields []
   From Field Email of Definition string: Parent 0xc000234c00 Field OF 0x0 Fields []
```

#### PrintFieldGraph
```
To Field Name of Definition string: Parent 0xc0003a3700 Field OF 0x0 Fields []
To Field UserID of Definition int: Parent 0xc0003a3800 Field OF 0xc00006a900 Fields []
From Field Name of Definition string: Parent 0xc0003a3900 Field OF 0x0 Fields []
From Field UserID of Definition int: Parent 0xc0003a3a00 Field OF 0x0 Fields []
To Field ID of Definition int: Parent 0xc00036e300 Field OF 0x0 Fields []
To Field Name of Definition string: Parent 0xc00036e400 Field OF 0x0 Fields []
To Field Email of Definition string: Parent 0xc000446200 Field OF 0x0 Fields []
From Field ID of Definition int: Parent 0xc000446300 Field OF 0x0 Fields []
From Field Name of Definition string: Parent 0xc000446400 Field OF 0x0 Fields []
From Field Email of Definition string: Parent 0xc000446500 Field OF 0x0 Fields []
```

#### PrintFieldTree

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

#### PrintFieldRelation

```
To Field int ID (domain.Account) and From Field int ID (models.Account) are related to each other.
To Field int ID (domain.Account) is not related to From Field string Name (models.Account).
To Field int ID (domain.Account) is not related to From Field string Email (models.Account).
To Field string Name (domain.Account) is not related to From Field int ID (models.Account).
To Field string Name (domain.Account) and From Field string Name (models.Account) are related to each other.
To Field string Name (domain.Account) is not related to From Field string Email (models.Account).
To Field string Email (domain.Account) is not related to From Field int ID (models.Account).
To Field string Email (domain.Account) is not related to From Field string Name (models.Account).
To Field string Email (domain.Account) and From Field string Email (models.Account) are related to each other.
To Field string Name (domain.Account) and From Field string Name (models.User) are related to each other.
To Field string Name (domain.Account) is not related to From Field int UserID (models.User).
To Field int UserID (domain.Account) is not related to From Field string Name (models.User).
To Field int UserID (domain.Account) and From Field int UserID (models.User) are related to each other.
```

#### CountFields

```
12
```