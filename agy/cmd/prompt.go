package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func Prompt(workspace string, args []string) {
	query := strings.Join(args, " ")
	if query == "" {
		query = "what should I work on?"
	}

	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║    AGY — Strategic Prompt            ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("  Query: %s\n\n", query)

	// Count project stats
	done := 0
	wip := 0
	planned := 0
	phases := map[string]int{}
	for _, feat := range f.Features {
		switch feat.Status {
		case plan.StatusDone:
			done++
		case plan.StatusWIP:
			wip++
		default:
			planned++
		}
		if feat.Phase != "" {
			phases[feat.Phase]++
		}
	}
	total := done + wip + planned

	fmt.Println("── Project Context ─────────────────────")
	fmt.Printf("  Total features: %d (%d done, %d in progress, %d planned)\n", total, done, wip, planned)
	fmt.Printf("  Phase distribution: ")
	for i := 1; i <= 8; i++ {
		key := fmt.Sprintf("Phase %d", i)
		if count, ok := phases[key]; ok && count > 0 {
			fmt.Printf("%s:%d ", key, count)
		}
	}
	fmt.Println()
	fmt.Println()

	// Suggest based on query
	queryLow := strings.ToLower(query)

	switch {
	case strings.Contains(queryLow, "next") || strings.Contains(queryLow, "work on") || strings.Contains(queryLow, "priority"):
		fmt.Println("── Suggested Next Work ─────────────────")
		fmt.Print(`  Based on the GSD-PLAN.md priority matrix, the highest-impact
  unstarted work is:

  1. BM25 keyword fallback (#25) — P0, Phase 2, 2 days
     Unlocks exact-match search. Critical for code/filesystem queries.
     Without this, searching "main()" returns bad results.

  2. Inotify file watcher (#15) — P0, Phase 2, 1 day
     Replaces 60s polling loop. Fixes the biggest UX complaint.

  3. ash workspace (#97) — P0, Phase 3, 3 days
     Core vibecoding UX. Makes the OS disposable.

  4. Cross-distro support (#57) — P0, Phase 4, 5 days
     Unlocks 80% of the Linux market. PM-translate maps.

  5. ash-ask RAG (#101) — P1, Phase 3, 3 days
     Killer feature. Natural language Q&A over your files.
`)
	case strings.Contains(queryLow, "risk") || strings.Contains(queryLow, "danger") || strings.Contains(queryLow, "blocker"):
		fmt.Println("── Risk Assessment ─────────────────────")
		fmt.Print(`  Current risks:
  • Cross-distro package drift — High likelihood, medium impact
  • Go rewrite stalls — Medium likelihood, high impact
  • Container mode perf gap — Medium likelihood, medium impact
  • LLM model deprecation — Medium likelihood, medium impact

  Mitigations in place:
  • CI tests all distros
  • Strangler fig pattern — bash stays as fallback
  • Multi-model support with --model flag
`)
	case strings.Contains(queryLow, "strateg") || strings.Contains(queryLow, "roadmap") || strings.Contains(queryLow, "long"):
		fmt.Println("── Strategic Roadmap ───────────────────")
		fmt.Print(`  Phase 1-3 (now): Install robustness, search quality, dev UX
  Phase 4-5   (next): Cross-platform, security hardening
  Phase 6-7   (later): Observability, federation
  Phase 8     (endgame): Plugin ecosystem, API surface

  Dependency chain: 1→2→3→8, 1→4→5→6→7
  Parallel tracks: 4 & 5 overlap with 3
`)
	default:
		fmt.Printf(`── Response ──────────────────────────────

  Project: Ash Linux — 106 features across 8 themes
  Progress: %d%% complete (%d/%d)

  %s

  Run 'agy status' for full project health.
  Run 'agy next' for priority-ordered build queue.
  Run 'agy contrast <dim>' to compare features.
  Run 'agy build <num>' for implementation plan.
`,
			done*100/total, done, total,
			suggestRelevantFeatures(queryLow, f),
		)
	}
}

func suggestRelevantFeatures(query string, p *plan.Plan) string {
	queryLow := strings.ToLower(query)
	matches := []string{}
	for _, feat := range p.Features {
		titleLow := strings.ToLower(feat.Title)
		if strings.Contains(titleLow, queryLow) || strings.Contains(queryLow, titleLow) {
			matches = append(matches, fmt.Sprintf("  • #%d %s", feat.Number, feat.Title))
		}
	}
	if len(matches) == 0 {
		return "  Ask about a specific topic (search, security, performance, etc.)\n  or run 'agy prompt \"what should I work on?\"' for task suggestions."
	}
	return strings.Join(matches, "\n")
}
