# Examples: Tests

Run unit and integration tests using `go test ./_tests` from `cd examples`.

## Integration Tests

The command line interface is straightforward. The loader uses a tested library. The matcher matches fields to other fields, which the generator depends on. Field-matching is heavily dependent on the `parser`, which provides the User Interface for end users _(developers)_. As a result, the `parser` contains the majority of edge cases this program encounters. Testing the entire program from end-to-end is more effective than unit tests _(with the exception of option-parsing)_.

| Test      | Description                                                          |
| :-------- | :------------------------------------------------------------------- |
| Alias     | Uses an alias import (for a copied struct).                          |
| Automap   | Uses the `automatch` option with a manual matcher option (`map`).    |
| Cast      | Uses the `cast` option for direct type conversion.                   |
| Cyclic    | Uses a nested struct (containing a field of the same type).          |
| Duplicate | Defines two structs with duplicate definitions, but not names.       |
| Import    | Imports a package in the setup file, that the output file exists in. |
| Multi     | Tests all types using multiple functions.                            |
| Option    | Tests Generator and Function option-parsing.                         |
| Same      | Generates an output file in the same directory as the setup file.    |

