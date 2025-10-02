package httpkit

import (
	"fmt"
	"net/http"
)


func WriteSSE(w http.ResponseWriter, data string) {
    fmt.Fprintf(w, "%s\n", data)
    if flusher, ok := w.(http.Flusher); ok {
        flusher.Flush()
    }
}

// SetSSEHeaders sets standard Server-Sent Event headers
func SetSSEHeaders(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
