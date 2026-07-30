package plan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type OllamaClient struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

func NewOllamaClient(baseURL string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   "llama3.2",
		Client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (o *OllamaClient) IsAvailable() bool {
	url := fmt.Sprintf("%s/api/tags", o.BaseURL)
	resp, err := o.Client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (o *OllamaClient) Generate(system, prompt string) (string, error) {
	req := OllamaRequest{
		Model:  o.Model,
		Prompt: prompt,
		System: system,
		Stream: false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", o.BaseURL)
	resp, err := o.Client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	var result OllamaResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}

func (o *OllamaClient) SuggestWork(p *Plan) string {
	if !o.IsAvailable() {
		return ""
	}

	systemPrompt := `You are a technical project manager analyzing a feature roadmap.
Given the project context, suggest the single most impactful feature to build next.
Be specific, concise, and explain WHY it's the priority.`

	done := 0
	for _, f := range p.Features {
		if f.Status == StatusDone {
			done++
		}
	}

	context := fmt.Sprintf(`Project: Ash Linux — Personal AI OS
Total features: %d (Done: %d, In progress: %d, Planned: %d)

P0 priorities:
- BM25 keyword fallback (#25) — 2 days, unlocks exact-match search
- Inotify file watcher (#15) — 1 day, fixes 60s poll latency
- ash workspace (#97) — 3 days, core vibecoding UX
- Cross-distro support (#57) — 5 days, unlocks 80%% of Linux market

Pick ONE feature to build next. Explain why in 2-3 sentences.`,
		len(p.Features), done, 0, len(p.Features)-done)

	response, err := o.Generate(systemPrompt, context)
	if err != nil {
		return ""
	}
	return response
}
