// Package pathkit provides filesystem path utilities.
//
// Includes functions for:
//   - Path existence and type checking (file, directory, symlink)
//   - Current working directory
//   - Path splitting and parsing
package pathkit

import (
	"os"
	"strings"
)


// CWD returns the current working directory.
func CWD() string {
    wd, err := os.Getwd()
    if err != nil {
        // fmt.Printf("Error getting working directory: %v\n", err)
        return ""
    }
    return wd
}


// Exists checks if a path exists.
func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}


// IsDir checks if path is a directory.
func IsDir(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}


// IsFile checks if path is a file.
func IsFile(path string) bool {
    info, err := os.Stat(path)
    return err == nil && !info.IsDir()
}


// IsSymlink checks if path is a symbolic link.
func IsSymlink(path string) bool {
    info, err := os.Lstat(path)
    return err == nil && info.Mode()&os.ModeSymlink != 0
}


// Split splits a path into segments, removing empty strings.
func Split(path string) []string {
    segments := make([]string, 0)
    for segment := range strings.SplitSeq(path, "/") {
        if segment == ""{
            continue
        }
        segments = append(segments, segment)
    }
    return segments
}

// QuickSplit splits a path into segments. Assumes path starts with "/".
func QuickSplit(path string) []string {
    segments := strings.Split(path, "/")[1:]
    last := len(segments) - 1
    if segments[last] == "" {
        segments = segments[:last]
    }

    return segments
}
