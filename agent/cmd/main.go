package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"vcx/agent/internal/config"
	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/db/dbsetup"
	"vcx/agent/internal/infra/fsmonitor"
	server "vcx/agent/internal/infra/http"
	"vcx/agent/internal/services/account"
	"vcx/agent/internal/services/migrations"
	"vcx/agent/internal/session"
	"vcx/pkg/logging"
)

func init() {
    logging.NewLogger(config.AppName)
}


func main() {
    // Ensure data directory exists before initializing DB
    if !dbsetup.PathExists() {
        log.Println("Created new data directory")
    }

    log.Printf("Initializing database at: %s", dbsetup.DefaultDBPath)
    db.Init(dbsetup.DefaultDBPath)

    // Run DB migrations
    if err := migrations.RunMigrations(context.Background()); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

	// Get or create default account
	baseCtx := context.Background()
	acc, err := account.GetOrCreateDefaultAccount(baseCtx)
	if err != nil {
		log.Fatalf("Failed to initialize account: %v", err)
	}

	// Add account to app context
	appCtx := session.WithAccountID(baseCtx, acc.ID)
	appCtx, appCancel := context.WithCancel(appCtx)
	var wg sync.WaitGroup

	startServer(appCtx, &wg)
	startMonitor(appCtx, &wg)

	waitForShutdown()

	log.Println("Shutting down...")
	appCancel()
	wg.Wait()
}


func startServer(appCtx context.Context, wg *sync.WaitGroup) {
    wg.Go(func() {
        server.Start(appCtx)
    })
}

func startMonitor(ctx context.Context, wg *sync.WaitGroup) {
	wg.Go(func() {
        fsmonitor.Start(ctx)
	})
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
