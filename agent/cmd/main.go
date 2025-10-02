package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"vcx/agent/internal/infra/db"
	"vcx/agent/internal/infra/fsmonitor"
	server "vcx/agent/internal/infra/http"
)

func init() {
    db.Init( "./journal.vcx?_journal_mode=WAL")
	// what is the state of this installation?
	// is there a db?
	// if not, initialize db
	// get relevant information into memory
	// - accountID/account object
	// - ? what else?
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	startServer(ctx, &wg)
	startMonitor(ctx, &wg)

	waitForShutdown()

	log.Println("Shutting down...")
	cancel()
	wg.Wait()
}


func startServer(ctx context.Context, wg *sync.WaitGroup) {
    wg.Go(func() {
        server.Start(ctx)
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
