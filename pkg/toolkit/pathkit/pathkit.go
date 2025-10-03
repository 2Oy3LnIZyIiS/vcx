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
