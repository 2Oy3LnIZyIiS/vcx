package project

import (
	"fmt"
	"net/http"
	"time"
	"vcx/pkg/toolkit/httpkit"
)

const APIPath = "/api/project"


func Handler() http.Handler {
    // Create submux for project routes
    mux := http.NewServeMux()

    // Register routes
    mux.HandleFunc("/init", initProject)
    mux.HandleFunc("/init-stream", initProjectStream)

    return http.StripPrefix(APIPath, mux)
}


func initProject(w http.ResponseWriter, r *http.Request) {
    httpkit.SetSSEHeaders(w)
    httpkit.WriteSSE(w, "init project called")


    // Simulate work delay
    time.Sleep(2 * time.Second)

    // Send completion event
    httpkit.WriteSSE(w, "{\"completed\": true}")
}



func initProjectStream(w http.ResponseWriter, r *http.Request) {
    httpkit.SetSSEHeaders(w)

    // Simulate project initialization steps
    steps := []string{
        "Creating project directory...",
        "Initializing version control...",
        "Setting up file monitoring...",
        "Creating initial snapshot...",
        "Project initialized successfully!",
    }

    for i, step := range steps {
        // Send progress update
        data := fmt.Sprintf("{\"step\": %d, \"total\": %d, \"message\": \"%s\"}", i+1, len(steps), step)
        httpkit.WriteSSE(w, data)

        // Simulate work delay
        time.Sleep(1 * time.Second)
    }

    // Send completion event
    httpkit.WriteSSE(w, "{\"completed\": true}")
}
