package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func Search(workspace string, args []string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	query := strings.Join(args, " ")
	if query == "" {
		fmt.Println("Usage: agy search <query>")
		fmt.Println("  Searches feature titles and descriptions for matching text.")
		return
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════╗\033[0m")
	fmt.Printf("\033[1;36m║  AGY — Search: %-28s║\033[0m\n", query)
	fmt.Println("\033[1;36m╚══════════════════════════════════════╝\033[0m")
	fmt.Println()

	queryLow := strings.ToLower(query)

	results := []plan.Feature{}
	seen := map[int]bool{}

	// Search by status keyword
	statusFilter := plan.StatusPlanned
	switch queryLow {
	case "done", "completed", "finished":
		statusFilter = plan.StatusDone
	case "wip", "progress", "working":
		statusFilter = plan.StatusWIP
	case "planned", "todo", "future":
		statusFilter = plan.StatusPlanned
	}

	for _, feat := range f.Features {
		match := false
		if strings.Contains(strings.ToLower(feat.Title), queryLow) {
			match = true
		}
		if strings.Contains(strings.ToLower(feat.Theme), queryLow) {
			match = true
		}
		// Match by status keyword
		if feat.Status == statusFilter && (queryLow == "done" || queryLow == "completed" || queryLow == "wip" || queryLow == "planned") {
			match = true
		}
		if match && !seen[feat.Number] {
			seen[feat.Number] = true
			results = append(results, feat)
		}
	}

	if len(results) == 0 {
		fmt.Println("  \033[1;30mNo features match your query.\033[0m")
		fmt.Println()
		fmt.Println("  Try: agy search workspace, agy search security, agy search gpu")
		return
	}

	fmt.Printf("  \033[1;37m%d matching features:\033[0m\n\n", len(results))

	for _, feat := range results {
		status := "\033[1;30m📅\033[0m"
		if feat.Status == plan.StatusDone {
			status = "\033[1;32m✅\033[0m"
		} else if feat.Status == plan.StatusWIP {
			status = "\033[1;33m🔧\033[0m"
		}

		// Highlight matching text
		title := feat.Title
		idx := strings.Index(strings.ToLower(title), queryLow)
		if idx >= 0 {
			title = title[:idx] + "\033[1;33m" + title[idx:idx+len(query)] + "\033[0m" + title[idx+len(query):]
		}

		fmt.Printf("  %s \033[1;37m#%d\033[0m %s\n", status, feat.Number, title)
	}
}
