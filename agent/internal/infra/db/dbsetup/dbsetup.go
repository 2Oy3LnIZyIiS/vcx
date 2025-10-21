package dbsetup

import (
	"os"
	"vcx/agent/internal/config"
	"vcx/agent/internal/infra/db/store/account"
	"vcx/agent/internal/infra/db/store/blob"
	"vcx/agent/internal/infra/db/store/branch"
	"vcx/agent/internal/infra/db/store/change"
	"vcx/agent/internal/infra/db/store/file"
	"vcx/agent/internal/infra/db/store/instance"
	"vcx/agent/internal/infra/db/store/project"
	"vcx/agent/internal/infra/db/store/simplekv"
	"vcx/agent/internal/infra/db/store/tag"
	"vcx/pkg/toolkit/pathkit"
)


var (
    DefaultDataPath = config.AppDataDir("data")
    DefaultDBPath   = config.AppDataDir("data", "journal.vcx") + "?_journal_mode=WAL"
)


func PathExists() bool {
    if !pathkit.Exists(DefaultDataPath) {
        os.MkdirAll(DefaultDataPath, os.ModePerm)
        return false
    }
    return true
}


func CreateTables() {
	account.CreateTable()
	blob.CreateTable()
	branch.CreateTable()
	change.CreateTable()
	file.CreateTable()
	instance.CreateTable()
	project.CreateTable()
	simplekv.CreateTable()
	tag.CreateTable()
}
