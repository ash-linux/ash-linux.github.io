package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)


func Build(workspace string, args []string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var targetNum int
	if len(args) > 0 {
		targetNum, _ = strconv.Atoi(args[0])
	}

	if targetNum == 0 {
		// Show all features that aren't done yet, grouped by priority
		fmt.Println("╔══════════════════════════════════════╗")
		fmt.Println("║    AGY — Build Queue                ║")
		fmt.Println("╚══════════════════════════════════════╝")
		fmt.Println()

		unbuilt := 0
		for _, feat := range f.Features {
			if feat.Status != plan.StatusDone {
				unbuilt++
			}
		}
		fmt.Printf("  %d features remaining to build\n\n", unbuilt)

		for _, feat := range f.Features {
			if feat.Status == plan.StatusDone {
				continue
			}
			phase := feat.Phase
			if phase == "" {
				phase = "TBD"
			}
			status := "📅"
			if feat.Status == plan.StatusWIP {
				status = "🔧"
			}
			fmt.Printf("  %s #%03d [%s] %s\n", status, feat.Number, phase, feat.Title)
		}
		return
	}

	// Show details for a specific feature
	var foundFeat *plan.Feature
	for i, feat := range f.Features {
		if feat.Number == targetNum {
			foundFeat = &f.Features[i]
			break
		}
	}
	if foundFeat == nil {
		fmt.Fprintf(os.Stderr, "Feature #%d not found\n", targetNum)
		os.Exit(1)
	}

	fmt.Printf("╔══════════════════════════════════════╗\n")
	fmt.Printf("║  Build Plan: #%d — %-30s║\n", foundFeat.Number, truncate(foundFeat.Title, 30))
	fmt.Printf("╚══════════════════════════════════════╝\n\n")

	statusStr := "📅 Planned"
	if foundFeat.Status == plan.StatusWIP {
		statusStr = "🔧 In Progress"
	} else if foundFeat.Status == plan.StatusDone {
		statusStr = "✅ Done"
	}
	fmt.Printf("  Status:  %s\n", statusStr)
	fmt.Printf("  Theme:   %s\n", foundFeat.Theme)
	fmt.Printf("  Phase:   %s\n", foundFeat.Phase)
	fmt.Println()

	if foundFeat.Status != plan.StatusDone {
		fmt.Println("── Suggested Implementation Steps ─────")
		fmt.Println(generateBuildSteps(foundFeat.Number))
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func generateBuildSteps(num int) string {
	steps := map[int]string{
		25: strings.TrimSpace(`
  1. Add SQLite FTS5 virtual table to LSFS daemon
  2. On file index, also insert title+content into FTS5
  3. Implement keyword_search() in lsfs-query
  4. Implement RRF fuse(semantic, keyword, top_k=20)
  5. Add --search-mode flag (hybrid|semantic|keyword)
  6. Test with exact-match queries (function names, paths)
`),
		15: strings.TrimSpace(`
  1. Install python-watchdog in ash-install.sh
  2. Replace polling loop with inotify observer
  3. Add 30s coalescing debounce for batch events
  4. Test with git clone (10k files) and incremental edits
`),
		97: strings.TrimSpace(`
  1. Create ash workspace create (btrfs snapshot of /home)
  2. Create ash workspace reset (restore snapshot)
  3. Create ash workspace list
  4. Wire into Super+Space as "workspace" mode
`),
		57: strings.TrimSpace(`
  1. Create pm_translate() in ash-install.sh: arch→deb→rpm→apk
  2. Write per-distro package map files in /etc/ash/pkgmaps/
  3. Update preflight to detect distro at runtime
  4. CI matrix: test install on Ubuntu, Fedora, Debian
`),
		101: strings.TrimSpace(`
  1. Create ash-ask CLI that takes a natural-language question
  2. Query Qdrant for top-10 relevant chunks
  3. Feed chunks + question to Ollama with QA prompt
  4. Return answer with source citations (filename + line)
`),
	}

	if s, ok := steps[num]; ok {
		return s
	}
	return strings.TrimSpace(`
  1. Define the feature's interface/side-effect contract
  2. Build the core logic (no dependencies on other Ash components)
  3. Add health check and error handling
  4. Wire into ash-install.sh or ship as standalone CLI
  5. Add to ash-doctor verification
  6. Update FUTURE.md status to ✅
`)
}
