package project

import (
	"fmt"
	"vcx/clients/cli/internal/client"
)

type Project struct {
	client.Client  // Embedded field
    Other string
}

func ProjectCommand(c client.Client) *Project {
	return &Project{
		Client: c,
	}
}


func (p *Project) Init() error {
	resp, err := p.HTTP.Get(p.BaseURL + "/api/project/init")
	if err != nil {
		return fmt.Errorf("error calling agent: %w", err)
	}
	defer resp.Body.Close()

    return client.HandleBody(resp, client.JustPrint)
}
