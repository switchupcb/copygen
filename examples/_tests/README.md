# Examples: Tests

The examples in this folder are used for testing.


| Test      | Description                                                       |
| :-------- | :---------------------------------------------------------------- |
| Alias     | Uses an alias import (for a copied struct).                       |
| Automap   | Uses the `automatch` option with a manual matcher option (`map`). |
| Cyclic    | Uses a nested struct (containing a field of the same type).       |
| Duplicate | Defines two structs with duplicate definitions, but not names.    |

## Integration Tests

The command line interface is straightforward. The loader uses a tested library. The matcher matches fields to other fields, which the generator depends on. Field-matching is heavily dependent on the `parser`, which provides the User Interface for end users _(developers)_. As a result, the `parser` contains the majority of edge cases this program encounters. Testing the entire program from end-to-end is more effective than unit tests for each package.

Run integration tests using `go test ./examples/_tests` from `/copygen`.