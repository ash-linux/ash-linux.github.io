package cmd

import (
	"fmt"
	"os"

	"github.com/anomalyco/ash-iso/agy/display"
	"github.com/anomalyco/ash-iso/agy/plan"
)

func Status(workspace string) {
	p, err := plan.LoadPlan(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║       AGY — Project Health           ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════╝\033[0m")
	fmt.Println()

	// Phase completion
	fmt.Println("\033[1;33m── Phase Completion ────────────────────\033[0m")
	completion := p.PhaseCompletion()
	phaseOrder := []string{"Phase 1", "Phase 2", "Phase 3", "Phase 4", "Phase 5", "Phase 6", "Phase 7", "Phase 8"}
	for _, name := range phaseOrder {
		if pct, ok := completion[name]; ok {
			bar := display.ProgressBar(pct, 20)
			fmt.Printf("  \033[1;37m%s:\033[0m %s\n", name, bar)
		}
	}
	fmt.Println()

	// Feature status overview
	fmt.Println("\033[1;33m── Feature Status (106 total) ─────────\033[0m")
	counts := f.StatusCounts()
	total := counts[plan.StatusDone] + counts[plan.StatusWIP] + counts[plan.StatusPlanned]
	fmt.Printf("  \033[1;32m✅ Done:\033[0m     %3d (%d%%)\n", counts[plan.StatusDone], counts[plan.StatusDone]*100/total)
	fmt.Printf("  \033[1;33m🔧 In Prog:\033[0m  %3d (%d%%)\n", counts[plan.StatusWIP], counts[plan.StatusWIP]*100/total)
	fmt.Printf("  \033[1;30m📅 Planned:\033[0m  %3d (%d%%)\n", counts[plan.StatusPlanned], counts[plan.StatusPlanned]*100/total)
	fmt.Println()

	// Phase deliverables detail
	for _, name := range phaseOrder {
		for _, phase := range p.Phases {
			if phase.Name == name {
				totalD := len(phase.Deliverables)
				doneD := 0
				for _, d := range phase.Deliverables {
					if d.Status == plan.StatusDone {
						doneD++
					}
				}
				color := "\033[1;32m"
				if doneD < totalD {
					color = "\033[1;33m"
				}
				if doneD == 0 {
					color = "\033[1;30m"
				}
				fmt.Printf("%s── %s (%d/%d done) ─────────────────\033[0m\n", color, name, doneD, totalD)
				for _, d := range phase.Deliverables {
					switch d.Status {
					case plan.StatusDone:
						fmt.Printf("  \033[1;32m✅\033[0m %s\n", d.Text)
					case plan.StatusWIP:
						fmt.Printf("  \033[1;33m🔧\033[0m %s\n", d.Text)
					default:
						fmt.Printf("    \033[1;30m%s\033[0m\n", d.Text)
					}
				}
				fmt.Println()
			}
		}
	}

	// Theme breakdown
	if len(f.FeatureThemes) > 0 {
		fmt.Println("\033[1;33m── Theme Breakdown ────────────────────\033[0m")
		t := display.Table{
			Headers: []string{"Theme", "Count"},
		}
		for _, theme := range f.FeatureThemes {
			count := 0
			for _, feat := range f.Features {
				if feat.Theme == theme {
					count++
				}
			}
			t.Rows = append(t.Rows, []string{theme, fmt.Sprintf("%d", count)})
		}
		fmt.Print(t.Render())
	}

	// Top 5 bang-for-buck
	fmt.Println()
	fmt.Println("\033[1;33m── Top 5 Bang-for-Buck ────────────────\033[0m")
	fmt.Println("  \033[1;37m1. #25\033[0m  BM25 keyword fallback")
	fmt.Println("  \033[1;37m2. #15\033[0m  Inotify file watcher")
	fmt.Println("  \033[1;37m3. #97\033[0m  ash workspace commands")
	fmt.Println("  \033[1;37m4. #57\033[0m  Cross-distro support")
	fmt.Println("  \033[1;37m5. #101\033[0m ash-ask RAG")
}
