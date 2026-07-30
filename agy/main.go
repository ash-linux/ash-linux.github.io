package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anomalyco/ash-iso/agy/cmd"
)

func main() {
	workspace := findWorkspace()
	if workspace == "" {
		fmt.Fprintln(os.Stderr, "\033[31mError:\033[0m cannot find project root (no GSD-PLAN.md or FUTURE.md found)")
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		cmd.Dashboard(workspace)
		return
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "status":
		cmd.Status(workspace)
	case "next":
		cmd.Next(workspace)
	case "contrast":
		cmd.Contrast(workspace, subArgs)
	case "build":
		cmd.Build(workspace, subArgs)
	case "prompt":
		cmd.Prompt(workspace, subArgs)
	case "dashboard", "dash":
		cmd.Dashboard(workspace)
	case "graph", "deps":
		cmd.Graph(workspace)
	case "timeline", "forecast":
		cmd.Timeline(workspace)
	case "search", "find":
		cmd.Search(workspace, subArgs)
	case "velocity", "git":
		cmd.GitVelocity(workspace)
	case "export", "json":
		handleExport(workspace, subArgs)
	case "csv":
		cmd.ExportCSV(workspace)
	case "markdown", "md":
		cmd.ExportMarkdown(workspace)
	case "install":
		cmd.SelfInstall()
	case "version", "--version", "-v":
		cmd.SelfVersion()
	case "sync":
		cmd.SyncStatus(workspace)
	case "help", "--help", "-h":
		printHelp()
	default:
		cmd.Prompt(workspace, args)
	}
}

func handleExport(workspace string, args []string) {
	if len(args) > 0 && args[0] == "--csv" {
		cmd.ExportCSV(workspace)
	} else if len(args) > 0 && args[0] == "--md" {
		cmd.ExportMarkdown(workspace)
	} else {
		cmd.ExportJSON(workspace, args)
	}
}

func findWorkspace() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		gsd := filepath.Join(dir, "GSD-PLAN.md")
		future := filepath.Join(dir, "FUTURE.md")
		if _, err := os.Stat(gsd); err == nil {
			if _, err := os.Stat(future); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func printHelp() {
	fmt.Println("\033[1;36magy\033[0m \033[1;30m— antigravity project monitor\033[0m")
	fmt.Println()
	fmt.Println("\033[1;33mMonitoring & Dashboard:\033[0m")
	fmt.Println("  \033[1;37magy\033[0m                    Project dashboard (default)")
	fmt.Println("  \033[1;37magy status\033[0m             Phase completion + feature status")
	fmt.Println("  \033[1;37magy dashboard\033[0m          Full project health overview")
	fmt.Println("  \033[1;37magy velocity\033[0m           Git activity & commit velocity")
	fmt.Println("  \033[1;37magy timeline\033[0m           Effort forecasting & timeline")
	fmt.Println()
	fmt.Println("\033[1;33mPlanning & Strategy:\033[0m")
	fmt.Println("  \033[1;37magy next\033[0m               Priority-ordered build queue")
	fmt.Println("  \033[1;37magy contrast [dim]\033[0m     Compare features by dimension")
	fmt.Println("  \033[1;37magy build <num>\033[0m        Implementation plan for feature")
	fmt.Println("  \033[1;37magy prompt <query>\033[0m     Strategic guidance (uses Ollama if available)")
	fmt.Println()
	fmt.Println("\033[1;33mAnalysis:\033[0m")
	fmt.Println("  \033[1;37magy graph\033[0m              Dependency graph visualization")
	fmt.Println("  \033[1;37magy search <query>\033[0m     Search features by keyword")
	fmt.Println()
	fmt.Println("\033[1;33mExport & Tools:\033[0m")
	fmt.Println("  \033[1;37magy json [--pretty]\033[0m    Export all data as JSON")
	fmt.Println("  \033[1;37magy csv\033[0m                Export features as CSV")
	fmt.Println("  \033[1;37magy markdown\033[0m            Export features as Markdown")
	fmt.Println("  \033[1;37magy sync\033[0m               Check FUTURE.md vs git history")
	fmt.Println("  \033[1;37magy install\033[0m            Install agy to system PATH")
	fmt.Println("  \033[1;37magy help\033[0m               This help")
	fmt.Println()
	fmt.Println("Dimensions for 'agy contrast': theme, status, phase, effort")
	fmt.Println()
	fmt.Println("agy reads GSD-PLAN.md + FUTURE.md from the project root.")
}

func init() {
	_ = strings.TrimSpace
}
