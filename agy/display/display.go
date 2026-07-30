package display

import (
	"fmt"
	"strings"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func (t Table) Render() string {
	if len(t.Rows) == 0 {
		return ""
	}
	cols := len(t.Headers)
	widths := make([]int, cols)
	for i, h := range t.Headers {
		widths[i] = len(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var b strings.Builder
	sep := func() {
		for i, w := range widths {
			if i > 0 {
				b.WriteString("├─")
			} else {
				b.WriteString("──")
			}
			b.WriteString(strings.Repeat("─", w))
			b.WriteString("─")
		}
		b.WriteString("──\n")
	}

	// header
	b.WriteString("──")
	for i, h := range t.Headers {
		if i > 0 {
			b.WriteString("─┬─")
		} else {
			b.WriteString("──")
		}
		b.WriteString(fmt.Sprintf("%-*s", widths[i], h))
	}
	b.WriteString("──\n")

	sep()

	for _, row := range t.Rows {
		b.WriteString("  ")
		for i, cell := range row {
			if i > 0 {
				b.WriteString(" │ ")
			}
			b.WriteString(fmt.Sprintf("%-*s", widths[i], cell))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func ProgressBar(pct float64, width int) string {
	filled := int(pct * float64(width) / 100)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("%s %3.0f%%", bar, pct)
}

func Bullet(status, text string) string {
	switch status {
	case "done":
		return fmt.Sprintf("  ✅ %s", text)
	case "wip":
		return fmt.Sprintf("  🔧 %s", text)
	default:
		return fmt.Sprintf("    %s", text)
	}
}
