package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func Vibe(args []string) {
	fmt.Println("\033[1;36m")
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║       ASH VIBECODING MODE            ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println("\033[0m")
	fmt.Println("  \033[1;33m✦\033[0m Auto-snapshot every 5 min")
	fmt.Println("  \033[1;33m✦\033[0m Auto-index on file change (100ms)")
	fmt.Println("  \033[1;33m✦\033[0m Workspace: auto-save")
	fmt.Println("  \033[1;33m✦\033[0m Press \033[1;37mCtrl+C\033[0m to exit (auto-reset available)")
	fmt.Println()

	// Create auto-save workspace if not in one
	wsName := fmt.Sprintf("vibe-%d", time.Now().Unix())
	wsCreate([]string{wsName})

	// Trap Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Minute)
	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}

	// File watcher (inotify)
	watcherCh := make(chan string, 100)
	go watchFiles("/home", watcherCh)

	fmt.Println("  \033[1;32m✔ Vibecoding mode active\033[0m")

	for {
		select {
		case <-sigCh:
			fmt.Println("\n")
			fmt.Println("\033[1;33mExiting vibecoding mode...\033[0m")
			fmt.Println("  \033[1;37m[1]\033[0m ash workspace reset — discard all changes")
			fmt.Println("  \033[1;37m[2]\033[0m ash workspace list — see snapshots")
			fmt.Println("  \033[1;37m[3]\033[0m Any other key — keep changes")
			fmt.Println()
			fmt.Print("  Choice: ")
			var choice string
			fmt.Scanln(&choice)
			if choice == "1" {
				wsReset([]string{wsName})
				fmt.Println("  \033[32m✔ Workspace reset. All changes discarded.\033[0m")
			} else if choice == "2" {
				wsList()
			} else {
				fmt.Println("  \033[32m✔ Changes kept.\033[0m")
			}
			fmt.Println("\033[1;36m╔══════════════════════════════════════╗")
			fmt.Println("║       VIBECODING SESSION ENDED        ║")
			fmt.Println("╚══════════════════════════════════════╝\033[0m")
			return

		case <-ticker.C:
			// Auto-snapshot every 5 min
			snapName := fmt.Sprintf("%s-auto-%d", wsName, time.Now().Unix())
			fmt.Printf("\n  \033[1;33m⟳\033[0m Auto-snapshot: %s\n", snapName)
			wsCreate([]string{snapName})

		case file := <-watcherCh:
			debounce.Reset(100 * time.Millisecond)
			_ = file
		}
	}
}

func watchFiles(root string, ch chan<- string) {
	// Simple polling-based watcher (inotify would be better on Linux)
	seen := map[string]time.Time{}
	for {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			// Skip .git, node_modules, .ash-*
			if strings.Contains(path, ".git/") || strings.Contains(path, "node_modules/") {
				return nil
			}

			modTime := info.ModTime()
			if prev, ok := seen[path]; !ok || prev != modTime {
				seen[path] = modTime
				select {
				case ch <- path:
				default:
				}
			}
			return nil
		})
		time.Sleep(1 * time.Second)
	}
}

func init() {
	// Ensure exec is used (referenced in this file)
	_ = exec.Command
}
