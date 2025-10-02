package main

import (
	"fmt"
	"os"
	"vcx/clients/cli/internal/client"
)

const agentURL = "http://localhost:9847"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vcx <command>")
		fmt.Println("Commands:")
		fmt.Println("  init         - Initialize project (simple)")
		fmt.Println("  init-progress - Initialize project (with progress)")
		os.Exit(1)
	}

	command := os.Args[1]
	cli := client.New(agentURL)

	switch command {
	case "init":
		callInit(cli)
	case "init-progress":
		callInitProgress(cli)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func callInit(cli *client.Client) {
	err := cli.Init()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func callInitProgress(cli *client.Client) {
	fmt.Println("Initializing project...")

	err := cli.InitWithProgress()
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		os.Exit(1)
	}
}
