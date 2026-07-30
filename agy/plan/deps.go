package plan

import "fmt"

// DepEdge represents a dependency between two features or phases.
type DepEdge struct {
	From string // depends on
	To   string // blocked by
}

// Dependencies defines the known dependency graph.
var Dependencies = []DepEdge{
	// Phase dependencies
	{"Phase 2", "Phase 1"},
	{"Phase 3", "Phase 2"},
	{"Phase 4", "Phase 1"},
	{"Phase 5", "Phase 4"},
	{"Phase 6", "Phase 5"},
	{"Phase 7", "Phase 6"},
	{"Phase 8", "Phase 3"},

	// Feature dependencies (specific features that block others)
	{"#25", "#15"},   // BM25 depends on inotify (needs fast reindex)
	{"#101", "#25"},  // ash-ask RAG depends on BM25
	{"#94", "#101"},  // Plugin system depends on ash-ask API
	{"#53", "#41"},   // ash-tui dashboard depends on install dashboard
	{"#90", "#85"},   // Federation depends on storage tiering
	{"#57", "#71"},   // Cross-distro depends on dependency solver
	{"#64", "#57"},   // Ansible export depends on cross-distro
	{"#97", "#21"},   // Workspace manager depends on swap-aware install
	{"#106", "#97"},  // VibeCoding depends on workspace manager
	{"#30", "#28"},   // File relationship graph depends on content extraction
	{"#31", "#15"},   // Clipboard history depends on inotify
	{"#62", "#57"},   // macOS depends on cross-distro patterns
	{"#63", "#62"},   // WSL2 depends on macOS patterns
	{"#77", "#73"},   // Verified boot depends on systemd hardening
	{"#78", "#77"},   // GPG signing depends on verified boot
	{"#79", "#74"},   // Secrets encryption depends on API auth
	{"#80", "#78"},   // USB signing depends on GPG signing
}

// BlockedBy returns which features block the given feature number.
func BlockedBy(num int) []string {
	prefix := fmt.Sprintf("#%d", num)
	var blockers []string
	for _, d := range Dependencies {
		if d.To == prefix {
			blockers = append(blockers, d.From)
		}
	}
	return blockers
}

// Blocks returns which features depend on the given feature number.
func Blocks(num int) []string {
	prefix := fmt.Sprintf("#%d", num)
	var dependents []string
	for _, d := range Dependencies {
		if d.From == prefix {
			dependents = append(dependents, d.To)
		}
	}
	return dependents
}
