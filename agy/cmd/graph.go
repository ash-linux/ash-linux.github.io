package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func Graph(workspace string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║       AGY — Dependency Graph                        ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()

	// Phase dependency chain
	fmt.Println("\033[1;33m── Phase Dependency Chain ───────────────────────────────\033[0m")
	fmt.Println()
	phases := []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Phase 5", "Phase 6", "Phase 7", "Phase 8"}

	// Level 1: 1 → 2 → 3 → 8 (core)
	fmt.Print("  ")
	for i, p := range []string{"Phase 1", "Phase 2", "Phase 3", "Phase 8"} {
		if i > 0 {
			fmt.Print(" \033[1;30m→\033[0m ")
		}
		fmt.Printf("\033[1;37m%s\033[0m", p)
	}
	fmt.Println(" \033[1;30m(Data platform)\033[0m")

	// Level 2: 1 → 4 → 5 → 6 → 7 (ops)
	fmt.Print("  ")
	for i, p := range []string{"Phase 1", "Phase 4", "Phase 5", "Phase 6", "Phase 7"} {
		if i > 0 {
			fmt.Print(" \033[1;30m→\033[0m ")
		}
		fmt.Printf("\033[1;37m%s\033[0m", p)
	}
	fmt.Println(" \033[1;30m(Platform)\033[0m")
	fmt.Println()

	// Feature dependency graph
	fmt.Println("\033[1;33m── Critical Feature Dependencies ────────────────────────\033[0m")
	fmt.Println()
	for _, dep := range plan.Dependencies {
		// Resolve to feature titles
		from := resolveTitle(dep.From, f)
		to := resolveTitle(dep.To, f)
		fmt.Printf("  \033[1;30m%s\033[0m \033[1;34m← requires ──\033[0m \033[1;37m%s\033[0m\n", dep.From+" "+from, dep.To+" "+to)
	}
	fmt.Println()

	// Blocked features summary
	fmt.Println("\033[1;33m── Bottleneck Analysis ──────────────────────────────────\033[0m")
	fmt.Println()
	bottleneckCount := map[string]int{}
	for _, dep := range plan.Dependencies {
		bottleneckCount[dep.From]++
	}
	for _, phase := range phases {
		found := false
		for _, feat := range f.Features {
			key := fmt.Sprintf("#%d", feat.Number)
			if c, ok := bottleneckCount[key]; ok && c >= 3 {
				if !found {
					fmt.Printf("  \033[1;31m%s\033[0m\n", phase)
					found = true
				}
				fmt.Printf("    \033[1;37m%s %s\033[0m — blocks %d other features\n", key, feat.Title, c)
			}
		}
	}
	if bottleneckCount["Phase 2"] > 0 {
		fmt.Println()
		fmt.Println("  \033[1;33mTip:\033[0m Build BM25 (#25) and inotify (#15) first to unblock everything.")
	}
}

func resolveTitle(key string, p *plan.Plan) string {
	var num int
	if _, err := fmt.Sscanf(key, "#%d", &num); err != nil {
		return ""
	}
	for _, feat := range p.Features {
		if feat.Number == num {
			title := feat.Title
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			return title
		}
	}
	return ""
}

func init() {
	_ = strings.TrimSpace
}
