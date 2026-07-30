package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anomalyco/ash-iso/agy/plan"
)

type ExportData struct {
	GeneratedAt string         `json:"generated_at"`
	Summary     ExportSummary  `json:"summary"`
	Phases      []ExportPhase  `json:"phases"`
	Features    []plan.Feature `json:"features"`
	Priorities  []plan.PriorityItem `json:"priorities"`
}

type ExportSummary struct {
	TotalFeatures   int            `json:"total_features"`
	Done            int            `json:"done"`
	InProgress      int            `json:"in_progress"`
	Planned         int            `json:"planned"`
	CompletionPct   float64        `json:"completion_pct"`
	PhaseCompletion map[string]float64 `json:"phase_completion"`
}

type ExportPhase struct {
	Name         string              `json:"name"`
	Total        int                 `json:"total"`
	Done         int                 `json:"done"`
	Deliverables []plan.Deliverable  `json:"deliverables"`
}

func ExportJSON(workspace string, args []string) {
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

	completion := p.PhaseCompletion()
	done := f.StatusCounts()[plan.StatusDone]
	wip := f.StatusCounts()[plan.StatusWIP]
	planned := f.StatusCounts()[plan.StatusPlanned]
	total := done + wip + planned

	var pct float64
	if total > 0 {
		pct = float64(done) / float64(total) * 100
	}

	var phases []ExportPhase
	for _, phase := range p.Phases {
		t := len(phase.Deliverables)
		d := 0
		for _, dl := range phase.Deliverables {
			if dl.Status == plan.StatusDone {
				d++
			}
		}
		phases = append(phases, ExportPhase{
			Name:         phase.Name,
			Total:        t,
			Done:         d,
			Deliverables: phase.Deliverables,
		})
	}

	data := ExportData{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Summary: ExportSummary{
			TotalFeatures:   total,
			Done:            done,
			InProgress:      wip,
			Planned:         planned,
			CompletionPct:   pct,
			PhaseCompletion: completion,
		},
		Phases:     phases,
		Features:   f.Features,
		Priorities: p.PriorityItems,
	}

	switch {
	case len(args) > 0 && args[0] == "--pretty":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(data)
	default:
		json.NewEncoder(os.Stdout).Encode(data)
	}
}

func ExportCSV(workspace string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("num,title,theme,status,phase")
	for _, feat := range f.Features {
		status := "planned"
		switch feat.Status {
		case plan.StatusDone:
			status = "done"
		case plan.StatusWIP:
			status = "wip"
		}
		title := strings.ReplaceAll(feat.Title, "\"", "'")
		theme := strings.ReplaceAll(feat.Theme, "\"", "")
		// Extract emoji + short name
		parts := strings.SplitN(theme, " ", 2)
		if len(parts) > 1 {
			theme = strings.TrimSpace(parts[1])
		}
		phase := feat.Phase
		if phase == "" {
			phase = "TBD"
		}
		fmt.Printf("%d,\"%s\",%s,%s,%s\n", feat.Number, title, theme, status, phase)
	}
}

func ExportMarkdown(workspace string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println("# Ash Linux — Feature Status")
	fmt.Println()
	fmt.Printf("_Generated: %s_\n\n", time.Now().Format("2006-01-02"))

	// Summary table
	fmt.Println("## Summary")
	fmt.Println()
	fmt.Println("| Status | Count |")
	fmt.Println("|--------|------:|")
	fmt.Printf("| ✅ Done | %d |\n", f.StatusCounts()[plan.StatusDone])
	fmt.Printf("| 🔧 In Progress | %d |\n", f.StatusCounts()[plan.StatusWIP])
	fmt.Printf("| 📅 Planned | %d |\n", f.StatusCounts()[plan.StatusPlanned])
	fmt.Println()

	// By theme
	fmt.Println("## By Theme")
	fmt.Println()
	fmt.Println("| Theme | Count |")
	fmt.Println("|-------|------:|")
	themeCount := map[string]int{}
	for _, feat := range f.Features {
		parts := strings.SplitN(feat.Theme, "(", 2)
		theme := strings.TrimSpace(parts[0])
		themeCount[theme]++
	}
	for theme, count := range themeCount {
		fmt.Printf("| %s | %d |\n", theme, count)
	}
	fmt.Println()

	// All features
	fmt.Println("## All Features")
	fmt.Println()
	fmt.Println("| # | Feature | Status | Phase |")
	fmt.Println("|---|---------|--------|-------|")
	for _, feat := range f.Features {
		status := "📅"
		switch feat.Status {
		case plan.StatusDone:
			status = "✅"
		case plan.StatusWIP:
			status = "🔧"
		}
		phase := feat.Phase
		if phase == "" {
			phase = "TBD"
		}
		fmt.Printf("| %d | %s | %s | %s |\n", feat.Number, feat.Title, status, phase)
	}
}
