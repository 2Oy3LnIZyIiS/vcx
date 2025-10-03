package project

import (
	"fmt"
	"net/http"
	"vcx/clients/cli/internal/client"
)



func Init() (*http.Response, error) {
    client := client.New()
    resp, err := client.Get("/api/project/init")
	if err != nil {
		return nil, fmt.Errorf("error calling agent: %w", err)
	}
    return resp, nil
}
