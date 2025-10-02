package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func registerRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/ping", ping)
    mux.HandleFunc("/health", health)
}

func health(w http.ResponseWriter, r *http.Request) {
    // Basic health check
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
}


func ping(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Only allow standard HTTP methods
    method := "UNKNOWN"
    switch r.Method {
    case "GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH", "TRACE", "CONNECT":
        method = r.Method
    }

    json.NewEncoder(w).Encode(map[string]string{
        "method": method,
        "timestamp": time.Now().Format(time.RFC3339Nano),
    })
}
