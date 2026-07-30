package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func SelfInstall() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	exe, err = filepath.Abs(exe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	targetDir := "/usr/local/bin"
	if runtime.GOOS == "darwin" {
		targetDir = "/opt/homebrew/bin"
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			targetDir = "/usr/local/bin"
		}
	}

	target := filepath.Join(targetDir, "agy")

	// Create symlink
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "\033[33mWarning: could not remove existing %s: %v\033[0m\n", target, err)
	}

	if err := os.Symlink(exe, target); err != nil {
		// Try copy instead
		cmd := exec.Command("cp", exe, target)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError installing: %v\n%s\033[0m\n", err, out)
			fmt.Println()
			fmt.Println("  Try: sudo agy install")
			os.Exit(1)
		}
	}

	fmt.Printf("\033[32m✔ agy installed to %s\033[0m\n", target)
	fmt.Println("  Run 'agy' from anywhere to monitor the project.")
}

func SelfVersion() {
	fmt.Println("agy — antigravity project monitor")
	fmt.Println("Version: 1.0.0 (built from ash-iso)")
	fmt.Println("See FUTURE.md feature #105")
}

func SyncStatus(workspace string) {
	// Read current FUTURE.md and check if we can auto-detect completed features
	// by checking git log for references
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	// Check git log for "#" feature references
	gitCmd := exec.Command("git", "log", "--oneline", "--format=%s", "-50")
	gitCmd.Dir = workspace
	out, err := gitCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[33mGit log not available\033[0m\n")
		return
	}

	commits := string(out)
	completed := 0
	for _, feat := range f.Features {
		if feat.Status == plan.StatusDone {
			continue
		}
		// Check if any commit message references this feature
		ref := fmt.Sprintf("#%d", feat.Number)
		if strings.Contains(commits, ref) {
			fmt.Printf("  \033[33m⚠ Feature #%d (%s) may be done but marked as planned\033[0m\n", feat.Number, feat.Title)
			completed++
		}
	}

	if completed == 0 {
		fmt.Println("  \033[32m✔ All features match git history\033[0m")
	} else {
		fmt.Printf("\n  \033[33m%d features may need status update in FUTURE.md\033[0m\n", completed)
	}
}
