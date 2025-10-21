package config

import (
	"path/filepath"
	"vcx/pkg/toolkit/systemkit"
)

const AppName = "vcx"

func AppDataDir(path ...string) string {
	return filepath.Join(append([]string{systemkit.DataDir(), AppName}, path...)...)
}