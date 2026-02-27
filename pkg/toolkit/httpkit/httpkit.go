// Package httpkit provides HTTP utilities for Server-Sent Events (SSE).
//
// Simplifies SSE implementation by handling:
//   - Standard SSE headers (Content-Type, Cache-Control, Connection)
//   - Message formatting and flushing
//   - CORS headers for cross-origin requests
package httpkit

import (
	"fmt"
	"net/http"
)


// WriteSSE writes a Server-Sent Event message and flushes the response.
func WriteSSE(w http.ResponseWriter, data string) {
    fmt.Fprintf(w, "data: %s\n\n", data)
    if flusher, ok := w.(http.Flusher); ok {
        flusher.Flush()
    }
}

// SetSSEHeaders sets standard Server-Sent Event headers.
func SetSSEHeaders(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
