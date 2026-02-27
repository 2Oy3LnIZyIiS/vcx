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
	"vcx/pkg/logging"
)

func init() {
    logging.NewLogger(config.AppName)
}


func main() {
    dbExists := dbsetup.PathExists()
    db.Init(dbsetup.DefaultDBPath)
    if !dbExists {
        dbsetup.CreateTables()
    }

	// what is the state of this installation?
	// is there a db?
	// if not, initialize db
	// get relevant information into memory
	// - accountID/account object
	// - ? what else?


    // TODO: Add accountID to appCtx

	appCtx, appCancel := context.WithCancel(context.Background())
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
