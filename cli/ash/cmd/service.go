package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type service struct {
	name string
	port string
	unit string
}

var services = []service{
	{name: "Qdrant", port: "6333", unit: "qdrant"},
	{name: "Ollama", port: "11434", unit: "ollama"},
	{name: "LSFS Daemon", port: "", unit: "lsfs-daemon"},
}

func Up(args []string) {
	profile := "desktop"
	if len(args) > 0 {
		profile = args[0]
	}

	fmt.Printf("Starting Ash services (%s profile)...\n", profile)

	// Filter services based on profile
	svcs := services
	if profile == "minimal" {
		svcs = services[:2] // Qdrant + Ollama only
	}

	for _, s := range svcs {
		fmt.Printf("  Starting %s... ", s.name)
		cmd := exec.Command("sudo", "systemctl", "start", s.unit)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("\033[31m✘\033[0m %v\n", strings.TrimSpace(string(out)))
		} else {
			fmt.Println("\033[32m✔\033[0m")
		}
	}

	// Wait for services to be healthy
	fmt.Println("\n  Waiting for services...")
	allHealthy := true
	for _, svc := range svcs {
		if svc.port == "" {
			fmt.Printf("  %-15s \033[33m⚠ no port check\033[0m\n", svc.name)
			continue
		}
		healthy := waitForPort(svc.port, 10*time.Second)
		if healthy {
			fmt.Printf("  %-15s \033[32m✔ healthy\033[0m\n", svc.name)
		} else {
			fmt.Printf("  %-15s \033[31m✘ unhealthy\033[0m\n", svc.name)
			allHealthy = false
		}
	}

	if allHealthy {
		fmt.Println("\n\033[32m✔ All Ash services are running\033[0m")
	} else {
		fmt.Println("\n\033[33m⚠ Some services may not be healthy. Run 'ash status' to check.\033[0m")
	}
}

func Down(args []string) {
	fmt.Println("Stopping Ash services...")

	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]
		fmt.Printf("  Stopping %s... ", s.name)
		cmd := exec.Command("sudo", "systemctl", "stop", s.unit)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("\033[31m✘\033[0m %v\n", strings.TrimSpace(string(out)))
		} else {
			fmt.Println("\033[32m✔\033[0m")
		}
	}
	fmt.Println("\033[32m✔ All Ash services stopped\033[0m")
}

func Status(args []string) {
	fmt.Println("Ash Service Status")
	fmt.Println("──────────────────")

	for _, s := range services {
		status := "\033[31m✘ stopped\033[0m"
		color := "\033[31m"

		cmd := exec.Command("systemctl", "is-active", s.unit)
		if out, err := cmd.Output(); err == nil {
			state := strings.TrimSpace(string(out))
			if state == "active" {
				status = "\033[32m✔ active\033[0m"
				color = "\033[32m"
			} else {
				status = fmt.Sprintf("\033[33m⚠ %s\033[0m", state)
				color = "\033[33m"
			}
		}

		// Also check port if available
		portStatus := ""
		if s.port != "" {
			if portReachable(s.port) {
				portStatus = "\033[32m●\033[0m"
			} else {
				portStatus = "\033[31m○\033[0m"
			}
		}

		fmt.Printf("  %s %-15s %s %s\n", portStatus, s.name, status, color)
	}
}

func waitForPort(port string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if portReachable(port) {
			return true
		}
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func portReachable(port string) bool {
	url := fmt.Sprintf("http://localhost:%s/health", port)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err == nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return true
	}
	return false
}
