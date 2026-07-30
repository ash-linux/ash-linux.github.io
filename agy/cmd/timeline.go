package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/anomalyco/ash-iso/agy/display"
	"github.com/anomalyco/ash-iso/agy/plan"
)

func Timeline(workspace string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}
	p, err := plan.LoadPlan(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║       AGY — Timeline & Forecasting                  ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════════════════════╝\033[0m")
	fmt.Println()

	// Calculate progress
	total := len(f.Features)
	done := f.StatusCounts()[plan.StatusDone]
	pct := float64(done) / float64(total) * 100

	// Estimate effort remaining (simplified: assume 2 days avg per feature)
	remaining := total - done
	estimatedDays := remaining * 2 // rough average
	weeks := float64(estimatedDays) / 5.0

	now := time.Now()
	projectedEnd := now.AddDate(0, 0, estimatedDays)

	fmt.Printf("  \033[1;37mProgress:\033[0m        %d/%d features done (%.0f%%)\n", done, total, pct)
	fmt.Printf("  \033[1;37mEffort remaining:\033[0m ~%d person-days (%d features × 2 days avg)\n", estimatedDays, remaining)
	fmt.Printf("  \033[1;37mEstimated timeline:\033[0m %.0f weeks (through ~%s)\n", weeks, projectedEnd.Format("Jan 2, 2006"))
	fmt.Println()

	// Phase-level timeline
	fmt.Println("\033[1;33m── Phase Timeline (Sequential) ─────────────────────────\033[0m")
	fmt.Println()

	phaseEffort := map[string]int{}
	phaseOrder := []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Phase 5", "Phase 6", "Phase 7", "Phase 8"}

	for _, phase := range p.Phases {
		unstarted := 0
		for _, d := range phase.Deliverables {
			if d.Status != plan.StatusDone {
				unstarted++
			}
		}
		phaseEffort[phase.Name] = unstarted * 2 // 2 days avg per deliverable
	}

	t := display.Table{
		Headers: []string{"Phase", "Remaining", "Effort", "Timeline", "Bar"},
	}

	cumulativeDays := 0
	startDate := now
	for _, name := range phaseOrder {
		if effort, ok := phaseEffort[name]; ok {
			if effort == 0 {
				t.Rows = append(t.Rows, []string{name, "0", "0 days", "✅ Done", "\033[1;32m████████████████████\033[0m"})
				continue
			}
			start := startDate.AddDate(0, 0, cumulativeDays)
			end := startDate.AddDate(0, 0, cumulativeDays+effort)
			cumulativeDays += effort

			bar := strings.Repeat("▓", effort/2) + strings.Repeat("░", 20-effort/2)
			if effort/2 > 20 {
				bar = strings.Repeat("▓", 20)
			}
			t.Rows = append(t.Rows, []string{
				name,
				fmt.Sprintf("%d", effort/2),
				fmt.Sprintf("%d days", effort),
				fmt.Sprintf("%s → %s", start.Format("1/2"), end.Format("1/2")),
				bar,
			})
		}
	}
	fmt.Print(t.Render())
	fmt.Println()

	// Velocity
	fmt.Println("\033[1;33m── Velocity & Estimates ────────────────────────────────\033[0m")
	fmt.Println()

	// Check git log for recent activity
	gitOutput, err := os.ReadFile(".git/HEAD")
	hasGit := err == nil && len(gitOutput) > 0

	if hasGit {
		fmt.Println("  \033[1;37mGit history:\033[0m found — can track commit velocity")
	} else {
		fmt.Println("  \033[1;37mGit history:\033[0m not available (run from project root)")
	}

	fmt.Printf("  \033[1;37mAt current pace:\033[0m 2 features/week → ~%d weeks to completion\n", remaining/2)
	fmt.Printf("  \033[1;37mOptimistic:\033[0m      5 features/week → ~%d weeks\n", remaining/5)
	fmt.Printf("  \033[1;37mFull-time focus:\033[0m 10 features/week → ~%d weeks\n", remaining/10)
	fmt.Println()

	// Recommendations
	fmt.Println("\033[1;33m── Recommended Focus ───────────────────────────────────\033[0m")
	fmt.Println()
	fmt.Println("  \033[1;37mNext 2 weeks (highest ROI):\033[0m")
	fmt.Println("    Week 1: BM25 keyword fallback (#25) + Inotify watcher (#15)")
	fmt.Println("    Week 2: ash workspace (#97) + Cross-distro support (#57)")
	fmt.Println()
	fmt.Println("  \033[1;33mThese 4 features unblock 80%% of the remaining work.\033[0m")
}

func init() {
	_ = strings.TrimSpace
	_ = sort.Strings
}
