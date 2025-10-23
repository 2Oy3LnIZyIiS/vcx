package debugkit

import (
	"fmt"
	"path/filepath"
	"runtime"
)


func GetCaller(skip int) string {
    _, file, line, ok := runtime.Caller(skip + 1)
    if !ok {
        return "unknown"
    }
    return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
