package server

import (
	"context"
	"log"
	"net/http"
	"time"
	"vcx/agent/internal/infra/http/api/project"
	"vcx/agent/internal/session"
)

func contextMiddleware(appCtx context.Context) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accountID, err := session.HasAccountID(appCtx)
			if err == nil {
				ctx = session.WithAccountID(ctx, accountID)
				r   = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Start(appCtx context.Context) {
	mux := http.NewServeMux()

	// Register routes
	registerRoutes(mux)
	mux.Handle(project.APIPath+"/", project.Handler())

	// Chain middleware
	handler := corsMiddleware(contextMiddleware(appCtx)(mux))

	server := &http.Server{
		Addr:    ":9847",
		Handler: handler,
	}

	go func() {
		log.Println("Server starting on :9847")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	<-appCtx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
}
