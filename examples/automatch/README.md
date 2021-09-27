# Example: Automatch

The automatch examples uses the automatcher to match three models with varying level of depth:
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

## YML

The YML specifies to use a depth of `1` for the domain Account. This means that the domain Account will have values set as far as a level `// 1` depth level in addition to its respective from fields. The YML specifies to use a depth of `0` for the models User (by default). This means that the domain User will only assign values (to the domain Account) if the field is at a depth level of `// 0`.

```yml
# Define where the code will be generated.
generated:
  filepath: ./copygen.go
  package: copygen

# Define the imports that are included in the generated file.
# Imports can also be defined in the type-property.
# import:
#  - github.com/switchupcb/copygen/examples/automatch/domain

# Define the functions to be generated.
# Properties with `# default` are NOT necessary to include (see Main).
functions:

 # Custom function options can be defined for template use (see Main).
 ModelsToDomain:

    # Define the types to be copied (to and from).
    # Custom type options (to and from) can be defined for template use (see Main).
    to:
      Account:
        # Define the import path for the type (required for automatch).
        import: github.com/switchupcb/copygen/examples/automatch/domain
        package:  domain      # default: none
        depth:    1           # default: 0

      
    from:
      User:
        import: github.com/switchupcb/copygen/examples/automatch/models
        package: models       # default: none
        depth:   2            # default: 0
        
      Account:
        import: github.com/switchupcb/copygen/examples/automatch/models
        package: models       # default: none
```

## Output

`copygen -yml path/to/yml`

```go
INSERT OUTPUT
```

## Output: Duplicate Fields

In the case that we set models User to a depth level of `// 2`, we would end up with duplicate fields.

```go
INSERT OUTPUT
```