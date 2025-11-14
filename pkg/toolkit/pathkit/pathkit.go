package pathkit

import (
	"os"
	"strings"
)


func CWD() string {
    wd, err := os.Getwd()
    if err != nil {
        // fmt.Printf("Error getting working directory: %v\n", err)
        return ""
    }
    return wd
}


// check if path exists
func Exists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}


// Check if path is a directory
func IsDir(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}


// Check if path is a file
func IsFile(path string) bool {
    info, err := os.Stat(path)
    return err == nil && !info.IsDir()
}


func IsSymlink(path string) bool {
    info, err := os.Lstat(path)
    return err == nil && info.Mode()&os.ModeSymlink != 0
}


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

// NOTE: Assumes path starts with a "/" and is a proper path
func QuickSplit(path string) []string {
    segments := strings.Split(path, "/")[1:]
    last := len(segments) - 1
    if segments[last] == "" {
        segments = segments[:last]
    }

    return segments
}
