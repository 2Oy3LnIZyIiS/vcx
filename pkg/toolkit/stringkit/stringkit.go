// Package stringkit provides string comparison utilities with case-insensitive options.
//
// Wraps standard string operations with optional case-insensitive matching.
// Useful for user input comparison and search functionality.
package stringkit

import "strings"


// Equals compares two strings with optional case-insensitive matching.
func Equals(a, b string, ignoreCase bool) bool {
    if ignoreCase {
        return strings.EqualFold(a, b)
    }
    return a == b
}


// Contains checks if string a contains string b with optional case-insensitive matching.
func Contains(a, b string, ignoreCase bool) bool {
    if ignoreCase {
        return strings.Contains(
            strings.ToLower(a),
            strings.ToLower(b),
        )
    }
    return strings.Contains(a, b)
}
