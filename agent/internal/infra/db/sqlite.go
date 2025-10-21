package db

import (
	// "database/sql"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"vcx/pkg/logging"
)

var (
    log *slog.Logger
    db  *sql.DB
)


func Init(dataSource string) {
    log = logging.GetLogger()
    var err error
    log.Debug(fmt.Sprintf("Attempting to open datasource at: %s", dataSource))
    db, err = sql.Open("sqlite3", dataSource)
    if err != nil {
        log.Error(fmt.Sprintf("Invalid DB config: %s", err))
        os.Exit(1)
    }

    if err = db.Ping(); err != nil {
        log.Error(fmt.Sprintf("DB unreachable: %s", err))
        os.Exit(1)
    }

    db.SetMaxOpenConns(1)  // SQLite only supports one writer at a time
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0) // connections are reused forever
}

func tableExists(tableName string) bool {
    if db == nil {
        return false
    }
    query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
    err := db.QueryRow(query, tableName).Scan(&tableName)
    return err == nil
}
