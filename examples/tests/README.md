# Examples: Tests

The examples in this folder are used for testing.


| Test   | Description                                                 |
| :----- | :---------------------------------------------------------- |
| Alias  | Uses an alias import (for a copied struct).                 |
| Cyclic | Uses a nested struct (containing a field of the same type). |

## Integration

The Command Line Interface and Generator (with templates) either work or don't work. The config uses a tested library. The matcher serves it purpose (matching fields to other fields) and is only effected by external options. For this reason, the majority of the edge cases this program encounters are located in the `parser`. In contrast, the parser is the most likely to change (as it is based on the UI). Integration tests are used in order to allow the parser to change without invalidating its tests.