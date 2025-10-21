package pathkit

import "os"


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
