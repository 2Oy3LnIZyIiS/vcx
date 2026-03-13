package dbsetup

import (
	"log"
	"os"
	"vcx/agent/internal/config"
	"vcx/pkg/toolkit/pathkit"
)

var (
	DefaultDataPath = config.AppDataDir("data")
	DefaultDBPath   = config.AppDataDir("data", "journal.vcx") + "?_journal_mode=WAL"
	BlobStorePath   = config.AppDataDir("data", "blobs")
)

func PathExists() bool {
    log.Printf("Checking data path: %s", DefaultDataPath)
    if !pathkit.Exists(DefaultDataPath) {
        log.Printf("Creating data directory: %s", DefaultDataPath)
        if err := os.MkdirAll(DefaultDataPath, os.ModePerm); err != nil {
            log.Printf("Failed to create data directory: %v", err)
            return false
        }
        return false
    }
    log.Printf("Data directory exists: %s", DefaultDataPath)
    return true
}
