package db

import (
	// "database/sql"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"vcx/pkg/logging"
)

var (
    log = logging.Log
    db  *sql.DB
)

func Init(dataSource string) {
    var err error
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
    query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
    var name string
    err := db.QueryRow(query, tableName).Scan(&name)
    return err == nil
}
