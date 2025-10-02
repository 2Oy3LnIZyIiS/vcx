// Package convert provides utility functions for converting between different data types.
// It aims to provide consistent, readable type conversions that can be used throughout
// an application to ensure uniform behavior.
package convertkit

// BoolToInt converts a boolean value to an integer.
// It returns 1 for true and 0 for false.
func BoolToInt(value bool) int {
    if value {
        return 1
    }
    return 0
}

// IntToBool converts an integer value to a boolean.
// It returns true for positive integers (> 0) and false for zero or negative integers.
func IntToBool(value int) bool {
    if value > 0 {
        return true
    }
    return false
}
