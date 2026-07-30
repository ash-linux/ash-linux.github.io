package plan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func LoadPlan(workspace string) (*Plan, error) {
	p := &Plan{}
	data, err := os.ReadFile(filepath.Join(workspace, "GSD-PLAN.md"))
	if err != nil {
		return nil, err
	}
	p.parsePhases(string(data))
	p.parsePriorities(string(data))
	return p, nil
}

var phaseHeaders = []string{
	"Phase 1", "Phase 2", "Phase 3", "Phase 4",
	"Phase 5", "Phase 6", "Phase 7", "Phase 8",
}

func (p *Plan) parsePhases(data string) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	var current *Phase
	deliverableRe := regexp.MustCompile(`^-\s*\[([ xX])\]\s*(.*)`)
	phaseHeaderRe := regexp.MustCompile(`^##\s+(Phase \d).*$`)

	for scanner.Scan() {
		line := scanner.Text()

		if matches := phaseHeaderRe.FindStringSubmatch(line); len(matches) > 1 {
			if current != nil {
				p.Phases = append(p.Phases, *current)
			}
			current = &Phase{Name: matches[1]}
			continue
		}

		if current != nil {
			if matches := deliverableRe.FindStringSubmatch(line); len(matches) > 2 {
				d := Deliverable{Text: strings.TrimSpace(matches[2])}
				ch := matches[1]
				if ch == "x" || ch == "X" {
					d.Status = StatusDone
				} else {
					d.Status = StatusPlanned
				}
				current.Deliverables = append(current.Deliverables, d)
			}

			if strings.Contains(line, "🔧") {
				if current.Deliverables != nil && len(current.Deliverables) > 0 {
					last := &current.Deliverables[len(current.Deliverables)-1]
					if last.Status != StatusDone {
						last.Status = StatusWIP
					}
				}
			}
		}
	}
	if current != nil {
		p.Phases = append(p.Phases, *current)
	}
}

func (p *Plan) parsePriorities(data string) {
	re := regexp.MustCompile(`\|[ \t]*(\w+\d*)[ \t]*\|[ \t]*(\d)[ \t]*\|[ \t]*([^|]*)[ \t]*\|[ \t]*([^|]*)[ \t]*\|[ \t]*([^|]*)[ \t]*\|`)
	lines := strings.Split(data, "\n")
	inTable := false
	for _, line := range lines {
		if strings.HasPrefix(line, "| Priority | Phase |") {
			inTable = true
			continue
		}
		if inTable && strings.HasPrefix(line, "|:") {
			continue
		}
		if inTable && strings.HasPrefix(line, "|") {
			matches := re.FindStringSubmatch(line)
			if len(matches) == 6 {
				item := PriorityItem{
					Priority: matches[1],
					Phase:    "Phase " + matches[2],
					Feature:  strings.TrimSpace(matches[3]),
					Effort:   strings.TrimSpace(matches[4]),
					Impact:   strings.TrimSpace(matches[5]),
				}
				p.PriorityItems = append(p.PriorityItems, item)
			}
		} else if inTable && !strings.HasPrefix(line, "|") {
			inTable = false
		}
	}
}
