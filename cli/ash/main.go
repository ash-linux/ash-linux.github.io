package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/ash-iso/cli/ash/cmd"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "workspace", "ws":
		cmd.Workspace(args[1:])
	case "up":
		cmd.Up(args[1:])
	case "down":
		cmd.Down(args[1:])
	case "status", "st":
		cmd.Status(args[1:])
	case "vibe":
		cmd.Vibe(args[1:])
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n\n", args[0])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("\033[1;36mash\033[0m \033[1;30m— Ash Linux CLI\033[0m")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  \033[1;37mash workspace create <name>\033[0m   Create a disposable dev environment")
	fmt.Println("  \033[1;37mash workspace reset <name>\033[0m    Roll back to clean state")
	fmt.Println("  \033[1;37mash workspace list\033[0m              List all workspaces")
	fmt.Println("  \033[1;37mash workspace delete <name>\033[0m    Delete a workspace")
	fmt.Println("  \033[1;37mash up\033[0m                         Start all Ash services")
	fmt.Println("  \033[1;37mash down\033[0m                       Stop all Ash services")
	fmt.Println("  \033[1;37mash status\033[0m                     Show service health")
	fmt.Println("  \033[1;37mash vibe\033[0m                       Enter vibecoding mode")
	fmt.Println()
	fmt.Println("\033[1;30mash is the unified CLI for Ash Linux — manage workspaces,")
	fmt.Println("orchestrate services, and enter vibecoding mode.")
}

func init() {
	_ = strings.TrimSpace
}
