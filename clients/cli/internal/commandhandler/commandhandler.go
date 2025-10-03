package commandhandler

import (
	"fmt"
	"os"
	"vcx/clients/cli/internal/client/api"
	"vcx/clients/cli/internal/client/api/project"
	"vcx/pkg/toolkit/pathkit"
)

func Route(args []string) {
	switch args[1] {
	case "init":
		Init(args)
	default:
		fmt.Printf("Unknown command: %s\n", args[1])
		os.Exit(1)
	}
}

func  Init(args []string) {
    cwd := pathkit.CWD()
    fmt.Println(cwd)

    resp, err := project.Init()
    if err != nil {
		fmt.Print(err)
		os.Exit(1)
    }
	defer resp.Body.Close()

    api.HandleBody(resp, api.JustPrint)
}
