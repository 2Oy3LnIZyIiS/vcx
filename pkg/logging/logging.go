package logging

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"vcx/pkg/toolkit/systemkit"
)


var ( addContext = true
      logLevel   *slog.LevelVar
      Log        *slog.Logger
      DEBUG      = slog.LevelDebug
      INFO       = slog.LevelInfo
      WARN       = slog.LevelWarn
      ERROR      = slog.LevelError
      FATAL      = slog.Level(12)
      Panic      = slog.Level(13)
      Fatal      = slog.Level(14) )


func initLogFile(name string) *os.File {
    os.MkdirAll( logPath(name), os.ModePerm )
    file, err := os.OpenFile( filepath.Join( logPath(name), name+".log" ),
							  os.O_CREATE|os.O_APPEND|os.O_WRONLY,
							  0644 )
    if err != nil {
        log.Fatal(err) }
    return file
}


func logPath(name string) string {
    return filepath.Join(systemkit.DataDir(), name, "logs")
}


func SetLogLevel(level slog.Level) {
    logLevel.Set(level)
}


func GetLogger() *slog.Logger {
    if Log == nil {
        return slog.Default()
    }
    return Log
}


func NewLogger(name string) *slog.Logger {
    logLevel = new(slog.LevelVar)  // .. Info by default
    Log = slog.New(
		slog.NewJSONHandler(
			initLogFile(name),
            &slog.HandlerOptions{ AddSource: addContext, Level: logLevel } ))
    slog.SetDefault(Log)
    SetLogLevel(DEBUG)

    return Log
}
