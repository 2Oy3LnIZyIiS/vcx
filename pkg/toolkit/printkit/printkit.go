// Package printkit provides terminal output utilities.
//
// ANSI escape sequence helpers for terminal manipulation.
// Currently supports clearing the current line.
package printkit


// ClearLine clears the current terminal line.
func ClearLine() {
    print("\r\033[K")
}
