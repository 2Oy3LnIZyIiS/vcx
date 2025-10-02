package server

import (
	"context"
	"log"
	"net/http"
	"time"
	"vcx/agent/internal/infra/http/api/project"
)

func Start(ctx context.Context) {
	mux := http.NewServeMux()

	// Register routes
	registerRoutes(mux)
    mux.Handle(project.APIPath+"/", project.Handler())

	// CORS for React client
    corsHandler := func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            return
        }
        mux.ServeHTTP(w, r)
    }

	server := &http.Server{
		Addr:    ":9847",
		Handler: http.HandlerFunc(corsHandler),
	}

	go func() {
		log.Println("Server starting on :9847")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}
