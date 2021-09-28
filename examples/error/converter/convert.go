// Package c provides custom converter functions used to copy types. The package is named in a special way for testing purposes.
package c

import (
	"strconv"
)

// Itoa calls the strconv.Itoa
func Itoa(i int) string {
	return strconv.Itoa(i)
}
