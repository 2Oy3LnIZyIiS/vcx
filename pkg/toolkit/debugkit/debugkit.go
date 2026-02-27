// Package debugkit provides runtime debugging utilities.
//
// Helps with debugging by providing caller information (file and line number).
// Useful for logging and error tracking.
package debugkit

import (
	"fmt"
	"path/filepath"
	"runtime"
)


// GetCaller returns the file and line number of the caller.
// skip specifies how many stack frames to skip (0 = immediate caller).
func GetCaller(skip int) string {
    _, file, line, ok := runtime.Caller(skip + 1)
    if !ok {
        return "unknown"
    }
    return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
