package systemkit

import (
	"go/build"
	"os"
	"path/filepath"
	"runtime"
)


func GetDefaultConcurrency() int {
    numCPU        := runtime.NumCPU()
    concurrent    := numCPU * 16
    maxConcurrent := 128
    if concurrent > maxConcurrent {
        concurrent = maxConcurrent
    }
    return concurrent
}


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
