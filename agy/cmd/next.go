package cmd

import (
	"fmt"
	"os"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func Next(workspace string) {
	p, err := plan.LoadPlan(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║       AGY — What To Build Next       ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════╝\033[0m")
	fmt.Println()

	// Priority table
	fmt.Println("\033[1;33m── Priority Sequencing ────────────────\033[0m")
	fmt.Printf("  \033[1;32mP0\033[0m \033[1;30mMust have — shipping this week\033[0m\n")
	fmt.Printf("  \033[1;33mP1\033[0m \033[1;30mImportant — next sprint\033[0m\n")
	fmt.Printf("  \033[1;30mP2\033[0m \033[1;30mNice to have — later\033[0m\n")
	fmt.Println()

	for _, item := range p.PriorityItems {
		pColor := "\033[1;30m"
		switch item.Priority {
		case "P0":
			pColor = "\033[1;32m"
		case "P1":
			pColor = "\033[1;33m"
		}
		fmt.Printf("  %s%s\033[0m [\033[1;37m%s\033[0m] \033[1;37m%s\033[0m — %s (%s)\n",
			pColor, item.Priority, item.Phase, item.Feature, item.Effort, item.Impact)
	}
	fmt.Println()

	// Suggested today's work
	fmt.Println("\033[1;33m── Suggested Today ────────────────────\033[0m")
	fmt.Println("  \033[1;37m1.\033[0m BM25 keyword fallback (\033[1;33m#25\033[0m) — \033[1;32m2 days\033[0m, unlocks exact-match")
	fmt.Println("  \033[1;37m2.\033[0m Inotify file watcher (\033[1;33m#15\033[0m)  — \033[1;32m1 day\033[0m, fixes 60s poll")
	fmt.Println("  \033[1;37m3.\033[0m ash workspace (\033[1;33m#97\033[0m)          — \033[1;32m3 days\033[0m, core vibecoding UX")
	fmt.Println()

	// P0 unstarted
	fmt.Println("\033[1;33m── P0 Items Not Yet Started ───────────\033[0m")
	for _, item := range p.PriorityItems {
		if item.Priority == "P0" {
			fmt.Printf("  \033[1;31m■\033[0m \033[1;37m%s\033[0m (\033[1;33m%s\033[0m) — %s\n", item.Feature, item.Effort, item.Impact)
		}
	}
}
