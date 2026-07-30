package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// ANSI Escape sequences
const (
	clearScreen = "\033[2J\033[H"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
	reset       = "\033[0m"
	bold        = "\033[1m"
	green       = "\033[32m"
	red         = "\033[31m"
	yellow      = "\033[33m"
	cyan        = "\033[36m"
)

func runCmd(command string) string {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func Dashboard(workspace string) {
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	// Handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print(showCursor)
		os.Exit(0)
	}()

	// Handle 'q' to quit - reading unbuffered is complex without termios,
	// so we use a simple goroutine reading from stty
	go func() {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		b := make([]byte, 1)
		for {
			os.Stdin.Read(b)
			if b[0] == 'q' || b[0] == 'Q' || b[0] == 3 {
				fmt.Print(showCursor)
				os.Exit(0)
			}
		}
	}()

	for {
		fmt.Print(clearScreen)
		
		fmt.Printf("%s%sAsh Linux - System Dashboard%s\n", bold, cyan, reset)
		fmt.Println(strings.Repeat("─", 50))
		
		// Services
		fmt.Printf("%sServices%s\n", bold, reset)
		
		qdrant := runCmd("curl -sf http://localhost:6333/healthz")
		if qdrant != "" {
			fmt.Printf("  %s●%s Qdrant  :6333 (Active)\n", green, reset)
		} else {
			fmt.Printf("  %s●%s Qdrant  :6333 (Down)\n", red, reset)
		}
		
		ollama := runCmd("curl -sf http://localhost:11434/api/version")
		if ollama != "" {
			fmt.Printf("  %s●%s Ollama  :11434 (Active)\n", green, reset)
		} else {
			fmt.Printf("  %s●%s Ollama  :11434 (Down)\n", red, reset)
		}
		
		lsfs := runCmd("systemctl --user is-active lsfs-daemon.service 2>/dev/null")
		if lsfs == "active" {
			fmt.Printf("  %s●%s LSFS Daemon (Active)\n", green, reset)
		} else {
			fmt.Printf("  %s●%s LSFS Daemon (Inactive)\n", yellow, reset)
		}
		fmt.Println()

		// Storage
		fmt.Printf("%sStorage%s\n", bold, reset)
		qSize := runCmd("du -sh /var/lib/qdrant 2>/dev/null | awk '{print $1}'")
		if qSize == "" {
			qSize = "Unknown"
		}
		indexedFiles := runCmd("curl -sf http://localhost:6333/collections/apps | grep -o '\"points_count\":[0-9]*' | cut -d: -f2")
		if indexedFiles == "" {
			indexedFiles = "0"
		}
		fmt.Printf("  Qdrant Size: %s\n", qSize)
		fmt.Printf("  Indexed Files: %s\n", indexedFiles)
		fmt.Printf("  Cache Size: %s\n", runCmd("du -sh ~/.cache/ash 2>/dev/null | awk '{print $1}'"))
		fmt.Println()

		// Model Info
		fmt.Printf("%sModel & Compute%s\n", bold, reset)
		model := runCmd("curl -sf http://localhost:11434/api/tags | grep -o '\"name\":\"[^\"]*' | head -1 | cut -d'\"' -f4")
		if model == "" {
			model = "None loaded"
		}
		fmt.Printf("  Current Model: %s\n", model)
		
		gpu := runCmd("lspci | grep -i 'vga\\|3d\\|display' | head -1 | cut -d':' -f3 | sed 's/^ *//'")
		if gpu == "" {
			gpu = "Unknown / CPU"
		}
		fmt.Printf("  GPU Status: %s\n", gpu)
		fmt.Println()

		// Recent Queries
		fmt.Printf("%sRecent Search Queries%s\n", bold, reset)
		queries := runCmd("tail -n 5 /var/log/ash/audit.log 2>/dev/null")
		if queries != "" {
			for _, line := range strings.Split(queries, "\n") {
				if line != "" {
					fmt.Printf("  > %s\n", line)
				}
			}
		} else {
			fmt.Println("  (No recent queries)")
		}
		
		fmt.Println()
		fmt.Printf("%s[ Press 'q' or Ctrl+C to exit ]%s\n", cyan, reset)
		
		time.Sleep(2 * time.Second)
	}
}
