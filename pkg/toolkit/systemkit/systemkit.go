// Package systemkit provides system information and cross-platform directory utilities.
//
// Handles platform differences for:
//   - Home directory (USERPROFILE on Windows, HOME on Unix)
//   - Application data directory (AppData/Roaming, .local/share, Library/Application Support)
//   - Optimal concurrency levels based on CPU count
package systemkit

import (
	"go/build"
	"os"
	"path/filepath"
	"runtime"
)


// GetDefaultConcurrency returns optimal concurrency level based on CPU count.
// Returns min(NumCPU * 16, 128).
func GetDefaultConcurrency() int {
    numCPU        := runtime.NumCPU()
    concurrent    := numCPU * 16
    maxConcurrent := 128
    if concurrent > maxConcurrent {
        concurrent = maxConcurrent
    }
    return concurrent
}


// HomeDir returns the user's home directory path.
func HomeDir() string {
	var path = ""
	if runtime.GOOS == "windows" {
		path = os.Getenv("USERPROFILE")
	} else if runtime.GOOS == "linux" {
		path = os.Getenv("HOME")
	} else if runtime.GOOS == "darwin" {
		path = os.Getenv("HOME")
	} else {
		path = build.Default.GOPATH
	}
	return path
}

// DataDir returns the platform-specific application data directory.
// Windows: AppData/Roaming, Linux: .local/share, macOS: Library/Application Support
func DataDir() string {
	var path = ""
	if runtime.GOOS == "windows" {
		path = "AppData/Roaming"
	} else if runtime.GOOS == "linux" {
		path = ".local/share"
	} else if runtime.GOOS == "darwin" {
		path = "Library/Application Support"
	} else {
		path = build.Default.GOPATH
	}

    return filepath.Join(HomeDir(), path)
}
