package plan

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func LoadFuture(workspace string) (*Plan, error) {
	p := &Plan{}
	data, err := os.ReadFile(filepath.Join(workspace, "FUTURE.md"))
	if err != nil {
		return nil, err
	}

	p.parseFeatures(string(data))
	p.parseThemes(string(data))
	return p, nil
}

func normalizeTheme(heading string) string {
	return strings.TrimSpace(strings.TrimPrefix(heading, "## "))
}

func (p *Plan) parseFeatures(data string) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	var currentTheme string

	featureRe := regexp.MustCompile(`^###\s+(\d+)\.\s+(.*)`)
	statusInHeadingRe := regexp.MustCompile(`- STATUS:\s*(✅|🔧|📅)\s*(.*?)$`)
	statusLineRe := regexp.MustCompile(`\*STATUS:\s*(.*?)\*`)
	phaseRe := regexp.MustCompile(`(Phase \d)`)
	for scanner.Scan() {
		line := scanner.Text()

		// Match theme headings (## but not ###)
		if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "### ") && !strings.HasPrefix(line, "## Theme") {
			currentTheme = normalizeTheme(line)
			continue
		}
		_ = currentTheme

		if matches := featureRe.FindStringSubmatch(line); len(matches) > 2 {
			num, _ := strconv.Atoi(matches[1])
			titleLine := matches[2]

			// Check for inline status in heading
			statusMatch := statusInHeadingRe.FindStringSubmatch(titleLine)
			title := titleLine
			status := StatusPlanned
			phase := ""

			if len(statusMatch) > 2 {
				emoji := statusMatch[1]
				statusText := strings.TrimSpace(statusMatch[2])
				title = strings.TrimSpace(strings.SplitN(titleLine, "- STATUS:", 2)[0])
				switch emoji {
				case "✅":
					status = StatusDone
				case "🔧":
					status = StatusWIP
					if pMatch := phaseRe.FindStringSubmatch(statusText); len(pMatch) > 1 {
						phase = pMatch[1]
					}
				case "📅":
					status = StatusPlanned
					if pMatch := phaseRe.FindStringSubmatch(statusText); len(pMatch) > 1 {
						phase = pMatch[1]
					}
				}
			}

			f := Feature{
				Number: num,
				Title:  title,
				Theme:  currentTheme,
				Status: status,
				Phase:  phase,
			}

			p.Features = append(p.Features, f)
			continue
		}

		// Also check separate status lines (backward compat)
		if len(p.Features) > 0 {
			last := &p.Features[len(p.Features)-1]
			if sMatch := statusLineRe.FindStringSubmatch(line); len(sMatch) > 1 && last.Status == StatusPlanned {
				statusStr := sMatch[1]
				if strings.HasPrefix(statusStr, "✅") {
					last.Status = StatusDone
				} else if strings.HasPrefix(statusStr, "🔧") {
					last.Status = StatusWIP
				}
			}
		}
	}
}

func (p *Plan) parseThemes(data string) {
	seen := map[string]bool{}
	for _, feat := range p.Features {
		if feat.Theme != "" && !seen[feat.Theme] {
			seen[feat.Theme] = true
			p.FeatureThemes = append(p.FeatureThemes, feat.Theme)
		}
	}
}

func (p *Plan) StatusCounts() map[Status]int {
	counts := map[Status]int{StatusDone: 0, StatusWIP: 0, StatusPlanned: 0}
	for _, f := range p.Features {
		counts[f.Status]++
	}
	return counts
}

func (p *Plan) PhaseCompletion() map[string]float64 {
	result := map[string]float64{}
	phaseCount := map[string]int{}
	phaseDone := map[string]int{}

	for _, phase := range p.Phases {
		phaseCount[phase.Name] = len(phase.Deliverables)
		done := 0
		for _, d := range phase.Deliverables {
			if d.Status == StatusDone {
				done++
			}
		}
		phaseDone[phase.Name] = done
	}

	for k, total := range phaseCount {
		if total > 0 {
			result[k] = float64(phaseDone[k]) / float64(total) * 100
		} else {
			result[k] = 0
		}
	}
	return result
}
