package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"vcx/pkg/toolkit/systemkit"
)


var ( addContext = true
      logLevel    *slog.LevelVar
      Log         *slog.Logger
      logFile     *os.File
      currentDate string
      appName     string
      DEBUG     = slog.LevelDebug
      INFO      = slog.LevelInfo
      WARN      = slog.LevelWarn
      ERROR     = slog.LevelError
      FATAL     = slog.Level(12)
      Panic     = slog.Level(13)
      Fatal     = slog.Level(14)
)


func getCurrentLogFileName(name string) string {
    return filepath.Join(logPath(name), fmt.Sprintf("%s.log", name))
}


func getArchivedLogFileName(name, date string) string {
    return filepath.Join(logPath(name), fmt.Sprintf("%s-%s.log", name, date))
}


func initLogFile(name string) *os.File {
    os.MkdirAll(logPath(name), os.ModePerm)
    file, err := os.OpenFile( getCurrentLogFileName(name),
                              os.O_CREATE|os.O_APPEND|os.O_WRONLY,
                              0644 )
    if err != nil {
        log.Fatal(err)
    }
    currentDate = time.Now().Format("2006-01-02")
    return file
}


func checkLogRotation() {
    today := time.Now().Format("2006-01-02")
    if today != currentDate && logFile != nil {
        // Close current file
        logFile.Close()

        // Rename current log file to archived name with yesterday's date
        os.Rename( getCurrentLogFileName(appName),
                   getArchivedLogFileName(appName, currentDate) )

        // Open new current log file
        logFile = initLogFile(appName)

        // Create new logger with new file
        Log = slog.New(
            slog.NewJSONHandler(
                logFile,
                &slog.HandlerOptions{AddSource: addContext, Level: logLevel}))
        slog.SetDefault(Log)
    }
}


func logPath(name string) string {
    return filepath.Join(systemkit.DataDir(), name, "logs")
}


func SetLogLevel(level slog.Level) {
    logLevel.Set(level)
}


func GetLogger() *slog.Logger {
    if Log == nil {
        // TODO: should call NewLogger with appropriate default name
        // ...   Maybe extract from go.mod adjust to local directory
        return slog.Default()
    }
    return Log
}


func NewLogger(name string) *slog.Logger {
    appName = name
    logLevel = new(slog.LevelVar)  // Info by default
    logFile = initLogFile(name)

    Log = slog.New(
        slog.NewJSONHandler(
            logFile,
            &slog.HandlerOptions{AddSource: addContext, Level: logLevel}))
    slog.SetDefault(Log)
    SetLogLevel(DEBUG)

    // Start daily rotation checker
    go func() {
        ticker := time.NewTicker(1 * time.Hour)  // Check every hour
        defer ticker.Stop()
        for range ticker.C {
            checkLogRotation()
        }
    }()

    return Log
}
