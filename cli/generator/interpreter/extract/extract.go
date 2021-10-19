// Package extract uses the `yaegi extract` tool in order to generate the reflect.Value symbols of internal types.
package extract

import "reflect"

// Symbols are extracted from the internal types (compiled at runtime).
var Symbols map[string]map[string]reflect.Value = make(map[string]map[string]reflect.Value)

//go:generate yaegi extract github.com/switchupcb/copygen/cli/models
