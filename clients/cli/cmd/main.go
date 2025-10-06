package main

import (
	"fmt"
	"os"
	"vcx/clients/cli/internal/commandhandler"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vcx <command>")
		fmt.Println("Commands:")
		fmt.Println("  init          - Initialize project (simple)")
		fmt.Println("  init-progress - Initialize project (with progress)")
		os.Exit(1)
	}

    commandhandler.Route(os.Args)
}
