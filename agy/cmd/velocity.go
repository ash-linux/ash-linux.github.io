package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type CommitInfo struct {
	Hash      string
	Date      time.Time
	Message   string
	Files     int
}

func GitVelocity(workspace string) {
	fmt.Println("\033[1;36m╔══════════════════════════════════════╗\033[0m")
	fmt.Println("\033[1;36m║  AGY — Git Velocity & Activity      ║\033[0m")
	fmt.Println("\033[1;36m╚══════════════════════════════════════╝\033[0m")
	fmt.Println()

	commits, err := getCommits(workspace)
	if err != nil {
		fmt.Printf("  \033[33mGit history not available: %v\033[0m\n", err)
		fmt.Println()
		fmt.Println("  Run this from the project root (a git repository).")
		return
	}

	if len(commits) == 0 {
		fmt.Println("  \033[33mNo commits found.\033[0m")
		return
	}

	// Calculate velocity
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	lastWeek := now.AddDate(0, 0, -7)

	monthly := 0
	weekly := 0
	filesChanged := 0

	for _, c := range commits {
		if c.Date.After(lastMonth) {
			monthly++
		}
		if c.Date.After(lastWeek) {
			weekly++
		}
		filesChanged += c.Files
	}

	// Calculate days since first commit
	firstCommit := commits[len(commits)-1].Date
	daysSince := now.Sub(firstCommit).Hours() / 24
	avgPerDay := float64(len(commits)) / daysSince

	fmt.Printf("  \033[1;37mTotal commits:\033[0m    %d\n", len(commits))
	fmt.Printf("  \033[1;37mFiles changed:\033[0m   %d\n", filesChanged)
	fmt.Printf("  \033[1;37mTimespan:\033[0m        %.0f days\n", daysSince)
	fmt.Printf("  \033[1;37mAve velocity:\033[0m    %.1f commits/day\n", avgPerDay)
	fmt.Println()
	fmt.Printf("  \033[1;37mLast 7 days:\033[0m     %d commits\n", weekly)
	fmt.Printf("  \033[1;37mLast 30 days:\033[0m    %d commits\n", monthly)
	fmt.Println()

	// Most active days
	fmt.Println("\033[1;33m── Recent Activity ──────────────────────\033[0m")
	fmt.Println()

	activityByDay := map[string]int{}
	for _, c := range commits {
		day := c.Date.Format("Mon Jan 2")
		activityByDay[day]++
	}

	type dayCount struct {
		day   string
		count int
	}
	var sortedDays []dayCount
	for d, c := range activityByDay {
		sortedDays = append(sortedDays, dayCount{d, c})
	}
	sort.Slice(sortedDays, func(i, j int) bool {
		return sortedDays[i].count > sortedDays[j].count
	})

	for i, dc := range sortedDays {
		if i >= 10 {
			break
		}
		n := dc.count
		if n < 0 {
			n = 0
		}
		if n > 20 {
			n = 20
		}
		bar := strings.Repeat("▓", n) + strings.Repeat("░", 20-n)
		fmt.Printf("  %-15s %s %d\n", dc.day, bar, dc.count)
	}
	fmt.Println()

	// Latest commits
	fmt.Println("\033[1;33m── Latest Commits ───────────────────────\033[0m")
	fmt.Println()
	limit := 8
	if len(commits) < limit {
		limit = len(commits)
	}
	for _, c := range commits[:limit] {
		msg := c.Message
		if len(msg) > 65 {
			msg = msg[:62] + "..."
		}
		h := c.Hash
		if len(h) > 8 {
			h = h[:8]
		}
		fmt.Printf("  \033[1;30m%s\033[0m \033[1;37m%s\033[0m\n", h, msg)
	}
}

func getCommits(workspace string) ([]CommitInfo, error) {
	cmd := exec.Command("git", "log", "--oneline", "--stat", "--date=iso-strict", "-100")
	cmd.Dir = workspace
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var commits []CommitInfo
	var current CommitInfo

	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			continue
		}
		if line == "" {
			continue
		}
		if len(line) > 0 && line[0] == ' ' {
			// stat line
			parts := strings.Fields(line)
			if len(parts) > 0 {
				current.Files++
			}
			continue
		}

		// New commit
		if current.Hash != "" {
			commits = append(commits, current)
		}
		current = CommitInfo{}

		// Parse "hash message" format (oneline)
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 1 {
			current.Hash = parts[0]
		}
		if len(parts) >= 2 {
			current.Message = strings.TrimSpace(parts[1])
		}
	}

	if current.Hash != "" {
		commits = append(commits, current)
	}

	// Get dates for commits — use short hashes to match
	dateCmd := exec.Command("git", "log", "--oneline", "--format=%h %ai", "-100")
	dateCmd.Dir = workspace
	dateOut, err := dateCmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(dateOut)), "\n")
		dateMap := map[string]time.Time{}
		for _, l := range lines {
			if len(l) < 10 {
				continue
			}
			spaceIdx := strings.Index(l, " ")
			if spaceIdx < 0 {
				continue
			}
			hash := l[:spaceIdx]
			dateStr := strings.TrimSpace(l[spaceIdx+1:])
			ts, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
			if err != nil {
				// Try without timezone
				if len(dateStr) >= 19 {
					ts, err = time.Parse("2006-01-02 15:04:05", dateStr[:19])
					if err != nil {
						continue
					}
				}
			}
			dateMap[hash] = ts
		}
		for i, c := range commits {
			if d, ok := dateMap[c.Hash]; ok {
				commits[i].Date = d
			}
		}
	}

	return commits, nil
}

func init() {
	_ = os.Args
}
