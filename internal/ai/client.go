package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/thenumber/app/internal/calc"
)

const SystemPrompt = `You are an educational allocation helper for Indian savers. You are not a SEBI-registered advisor and must not give personalized financial advice.

Return ONLY JSON of this shape:
{"sleeves":[{"category":"Large-cap index","percent":60}],"rationale":"..."}

Rules:
- Use only these category names: Large-cap index, Mid-cap index, Small-cap index, Debt / government securities, Gold, International equity index, Cash / liquid.
- Percents must be numbers that sum to 100.
- Never name stocks, mutual funds, ETFs, AMCs, tickers, or specific products.
- Keep rationale short, category-level, and educational.`

type GuideInput struct {
	Age               int
	Horizon           int
	Risk              string
	MonthlyInvestable float64
	Goals             string
}

type GuideResult struct {
	Sleeves   []calc.Sleeve `json:"sleeves"`
	Rationale string        `json:"rationale"`
	Source    string        `json:"source"`
	Warning   string        `json:"warning,omitempty"`
}

func Guide(ctx context.Context, in GuideInput) GuideResult {
	fallback := calc.HeuristicAllocation(in.Age, in.Horizon, in.Risk)
	base := GuideResult{
		Sleeves:   fallback.Sleeves,
		Rationale: fallback.Rationale,
		Source:    "On-device heuristic",
	}

	client := &http.Client{Timeout: 20 * time.Second}

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		model := getenv("OPENAI_MODEL", "gpt-4o-mini")
		parsed, err := callOpenAI(ctx, client, key, model, in)
		if err == nil {
			return GuideResult{Sleeves: parsed.Sleeves, Rationale: parsed.Rationale, Source: "Model"}
		}
		base.Warning = "Could not reach the model. Showing the on-device heuristic instead."
		return base
	}

	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		model := getenv("ANTHROPIC_MODEL", "claude-3-5-haiku-latest")
		parsed, err := callAnthropic(ctx, client, key, model, in)
		if err == nil {
			return GuideResult{Sleeves: parsed.Sleeves, Rationale: parsed.Rationale, Source: "Model"}
		}
		base.Warning = "Could not reach the model. Showing the on-device heuristic instead."
		return base
	}

	return base
}

func userPrompt(in GuideInput) string {
	return fmt.Sprintf(
		"Age: %d\nInvesting horizon (years): %d\nRisk stance: %s\nMonthly investable (INR): %.0f\nGoals (user text): %s\n\nPropose a category-level mix. No product names.",
		in.Age, in.Horizon, in.Risk, in.MonthlyInvestable, in.Goals,
	)
}

func callOpenAI(ctx context.Context, client *http.Client, key, model string, in GuideInput) (GuidanceJSON, error) {
	body := map[string]any{
		"model":           model,
		"temperature":     0.3,
		"response_format": map[string]string{"type": "json_object"},
		"messages": []map[string]string{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": userPrompt(in)},
		},
	}
	raw, err := postJSON(ctx, client, "https://api.openai.com/v1/chat/completions", body, map[string]string{
		"Authorization": "Bearer " + key,
	})
	if err != nil {
		return GuidanceJSON{}, err
	}
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return GuidanceJSON{}, err
	}
	if len(resp.Choices) == 0 {
		return GuidanceJSON{}, errString("empty OpenAI response")
	}
	return ParseGuidanceJSON(resp.Choices[0].Message.Content)
}

func callAnthropic(ctx context.Context, client *http.Client, key, model string, in GuideInput) (GuidanceJSON, error) {
	body := map[string]any{
		"model":       model,
		"max_tokens":  800,
		"temperature": 0.3,
		"system":      SystemPrompt,
		"messages": []map[string]any{
			{"role": "user", "content": userPrompt(in)},
		},
	}
	raw, err := postJSON(ctx, client, "https://api.anthropic.com/v1/messages", body, map[string]string{
		"x-api-key":         key,
		"anthropic-version": "2023-06-01",
	})
	if err != nil {
		return GuidanceJSON{}, err
	}
	var resp struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return GuidanceJSON{}, err
	}
	text := ""
	for _, c := range resp.Content {
		if c.Type == "text" {
			text += c.Text
		}
	}
	if text == "" {
		return GuidanceJSON{}, errString("empty Anthropic response")
	}
	return ParseGuidanceJSON(text)
}

func postJSON(ctx context.Context, client *http.Client, url string, body any, headers map[string]string) ([]byte, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("api status %d", res.StatusCode)
	}
	return raw, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
