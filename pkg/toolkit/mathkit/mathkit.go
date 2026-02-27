// Package mathkit provides basic math utilities.
//
// Simple math operations not included in the standard library.
package mathkit


// Abs returns the absolute value of an integer.
func Abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
