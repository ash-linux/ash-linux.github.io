package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const workspacesDir = "/home/.ash/workspaces"

func Workspace(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ash workspace <create|reset|list|delete> [name]")
		return
	}

	switch args[0] {
	case "create":
		wsCreate(args[1:])
	case "reset":
		wsReset(args[1:])
	case "list":
		wsList()
	case "delete":
		wsDelete(args[1:])
	default:
		fmt.Printf("Unknown workspace command: %s\n", args[0])
	}
}

func wsCreate(args []string) {
	name := "default"
	if len(args) > 0 {
		name = args[0]
	}

	fmt.Printf("Creating workspace '%s'...\n", name)

	// Check if btrfs is available
	if !commandExists("btrfs") {
		fmt.Println("\033[33m⚠ btrfs not found — using directory snapshot (copy-on-write not available)\033[0m")
		wsCreateFallback(name)
		return
	}

	// Create workspace directory
	wsDir := filepath.Join(workspacesDir, name)
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating workspace dir: %v\n", err)
		os.Exit(1)
	}

	// Create btrfs snapshot of /home
	homeSnap := filepath.Join(wsDir, "home.snap")
	cmd := exec.Command("btrfs", "subvolume", "snapshot", "-r", "/home", homeSnap)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "btrfs snapshot failed: %v\n%s\n", err, out)
		os.Exit(1)
	}

	// Store metadata
	metaFile := filepath.Join(wsDir, ".ash-workspace")
	meta := fmt.Sprintf("name=%s\ncreated=%s\nsubvol=%s\n", name, time.Now().Format(time.RFC3339), homeSnap)
	os.WriteFile(metaFile, []byte(meta), 0644)

	fmt.Printf("\033[32m✔ Workspace '%s' created\033[0m\n", name)
	fmt.Printf("  Snapshot: %s\n", homeSnap)
	fmt.Println("  Run 'ash workspace reset', 'ash workspace delete'")
}

func wsCreateFallback(name string) {
	// Fallback: copy /home to a snapshot dir using cp --reflink if available
	wsDir := filepath.Join("/tmp/.ash-workspaces", name)
	homeSnap := filepath.Join(wsDir, "home.snap")

	if err := os.MkdirAll(wsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("cp", "-a", "--reflink=auto", "/home", homeSnap)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "snapshot failed: %v\n%s\n", err, out)
		os.Exit(1)
	}

	metaFile := filepath.Join(wsDir, ".ash-workspace")
	meta := fmt.Sprintf("name=%s\ncreated=%s\nfallback=%s\n", name, time.Now().Format(time.RFC3339), homeSnap)
	os.WriteFile(metaFile, []byte(meta), 0644)

	fmt.Printf("\033[32m✔ Workspace '%s' created (fallback mode)\033[0m\n", name)
}

func wsReset(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ash workspace reset <name>")
		return
	}
	name := args[0]
	fmt.Printf("Resetting workspace '%s'...\n", name)

	wsDir := locateWorkspace(name)
	if wsDir == "" {
		fmt.Fprintf(os.Stderr, "Workspace '%s' not found\n", name)
		os.Exit(1)
	}

	metaFile := filepath.Join(wsDir, ".ash-workspace")
	meta, err := os.ReadFile(metaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read workspace metadata: %v\n", err)
		os.Exit(1)
	}

	if strings.Contains(string(meta), "subvol=") {
		// btrfs mode: restore snapshot
		homeSnap := filepath.Join(wsDir, "home.snap")
		cmd := exec.Command("btrfs", "subvolume", "snapshot", homeSnap, "/home")
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "restore failed: %v\n%s\n", err, out)
			os.Exit(1)
		}
	} else if strings.Contains(string(meta), "fallback=") {
		homeSnap := filepath.Join(wsDir, "home.snap")
		cmd := exec.Command("rsync", "-a", "--delete", homeSnap+"/", "/home/")
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "restore failed: %v\n%s\n", err, out)
			os.Exit(1)
		}
	}

	fmt.Printf("\033[32m✔ Workspace '%s' reset to clean state\033[0m\n", name)
}

func wsList() {
	dirs := []string{workspacesDir, "/tmp/.ash-workspaces"}

	found := false
	for _, d := range dirs {
		entries, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			metaFile := filepath.Join(d, e.Name(), ".ash-workspace")
			meta, err := os.ReadFile(metaFile)
			if err != nil {
				continue
			}
			found = true
			mode := "btrfs"
			if strings.Contains(string(meta), "fallback") {
				mode = "copy"
			}
			var created string
			if _, err := fmt.Sscanf(string(meta), "name=%*s\ncreated=%s", &created); err == nil {
				created = created[:19] // truncate to date
			} else {
				created = "unknown"
			}
			fmt.Printf("  \033[1;37m%-20s\033[0m %s  [%s]\n", e.Name(), created, mode)
		}
	}

	if !found {
		fmt.Println("  No workspaces found. Create one with 'ash workspace create <name>'")
	}
}

func wsDelete(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ash workspace delete <name>")
		return
	}
	name := args[0]

	wsDir := locateWorkspace(name)
	if wsDir == "" {
		fmt.Fprintf(os.Stderr, "Workspace '%s' not found\n", name)
		return
	}

	meta, _ := os.ReadFile(filepath.Join(wsDir, ".ash-workspace"))
	if strings.Contains(string(meta), "subvol=") {
		homeSnap := filepath.Join(wsDir, "home.snap")
		exec.Command("btrfs", "subvolume", "delete", homeSnap).Run()
	}

	os.RemoveAll(wsDir)
	fmt.Printf("\033[32m✔ Workspace '%s' deleted\033[0m\n", name)
}

func locateWorkspace(name string) string {
	dirs := []string{workspacesDir, "/tmp/.ash-workspaces"}
	for _, d := range dirs {
		wsDir := filepath.Join(d, name)
		if _, err := os.Stat(filepath.Join(wsDir, ".ash-workspace")); err == nil {
			return wsDir
		}
	}
	return ""
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
