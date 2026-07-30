package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anomalyco/ash-iso/agy/plan"
)

func Contrast(workspace string, args []string) {
	f, err := plan.LoadFuture(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
		os.Exit(1)
	}

	dimension := "theme"
	if len(args) > 0 {
		dimension = args[0]
	}

	fmt.Println("\033[1;36m╔══════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║    AGY — Feature Contrast Matrix     ║\033[0m")
	fmt.Printf("\033[1;36m║    Dimension: %-22s║\033[0m\n", dimension)
	fmt.Println("\033[1;36m╚══════════════════════════════════════╝\033[0m")
	fmt.Println()

	groupFn := getGroupFn(dimension)

	groups := map[string][]plan.Feature{}
	for _, feat := range f.Features {
		key := groupFn(feat)
		groups[key] = append(groups[key], feat)
	}

	keys := sortedKeys(groups)
	for _, key := range keys {
		feats := groups[key]
		fmt.Printf("\033[1;33m── %s (%d) ────────────────────────\033[0m\n", key, len(feats))
		for _, feat := range feats {
			status := "\033[1;30m📅\033[0m"
			if feat.Status == plan.StatusDone {
				status = "\033[1;32m✅\033[0m"
			} else if feat.Status == plan.StatusWIP {
				status = "\033[1;33m🔧\033[0m"
			}
			fmt.Printf("  %s \033[1;37m#%d\033[0m %s\n", status, feat.Number, feat.Title)
		}
		fmt.Println()
	}
}

func getGroupFn(dimension string) func(plan.Feature) string {
	switch dimension {
	case "status":
		return func(f plan.Feature) string {
			switch f.Status {
			case plan.StatusDone:
				return "✅ Done"
			case plan.StatusWIP:
				return "🔧 In Progress"
			default:
				return "📅 Planned"
			}
		}
	case "phase":
		return func(f plan.Feature) string {
			if f.Phase != "" {
				return f.Phase
			}
			return "TBD"
		}
	case "effort":
		return func(f plan.Feature) string {
			n := f.Number
			switch {
			case n <= 12:
				return "1-3 days"
			case n <= 24:
				return "2-5 days"
			case n <= 40:
				return "3-7 days"
			case n <= 56:
				return "1-3 days"
			case n <= 72:
				return "3-10 days"
			case n <= 84:
				return "2-5 days"
			case n <= 93:
				return "2-5 days"
			default:
				return "3-7 days"
			}
		}
	default: // theme
		return func(f plan.Feature) string {
			parts := strings.SplitN(f.Theme, "(", 2)
			return strings.TrimSpace(parts[0])
		}
	}
}

func sortedKeys(m map[string][]plan.Feature) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort: put "TBD" last, done first
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func init() {
	_ = os.Args
}
