package reasoning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"rationalgo/internal/models"
	"rationalgo/internal/services/decision"
)

const anthropicAPI = "https://api.anthropic.com/v1/messages"

// Service calls the Anthropic API to produce a structured DecisionRecord.
type Service struct {
	APIKey string
	Model  string
	Client *http.Client
}

// New returns a reasoning Service configured with the given API key.
func New(apiKey string) *Service {
	return &Service{
		APIKey: apiKey,
		Model:  "claude-sonnet-4-6",
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

// ── Anthropic wire types ──────────────────────────────────────────────────────

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// anthropicResponse is the top-level shape the API always returns.
type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"` // "text" for normal completions
		Text string `json:"text"`
	} `json:"content"`
	// Error is non-nil when the API returns a 4xx/5xx with a JSON error body.
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ── What we ask the model to return ──────────────────────────────────────────
// Separate from DecisionRecord so the model only fills in the reasoning fields.
// Status, Policy, Timestamp, ReasoningHash are set by code after parsing.

type reasoningOutput struct {
	VendorChosen  models.VendorOption  `json:"VendorChosen"`
	Alternatives  []models.Alternative `json:"Alternatives"`
	ExpectedValue string               `json:"ExpectedValue"`
	Confidence    float64              `json:"Confidence"`
}

// ── Public API ────────────────────────────────────────────────────────────────

// GenerateDecision calls the LLM, parses its JSON output, then assembles a
// full DecisionRecord by merging reasoning output with policy and computed fields.
func (s *Service) GenerateDecision(
	ctx context.Context,
	agentID string,
	sessionID string,
	intent string,
	vendors []models.VendorOption,
	policy models.PolicyResult,
) (*models.DecisionRecord, error) {
	if s.APIKey == "" {
		return nil, fmt.Errorf("reasoning: RATIONALGO_ANTHROPIC_KEY is not set")
	}

	prompt := buildPrompt(intent, vendors, policy)

	raw, err := s.call(ctx, prompt)
	if err != nil {
		return nil, err
	}

	out, err := parseOutput(raw)
	if err != nil {
		return nil, err
	}

	return assembleRecord(agentID, sessionID, intent, out, policy), nil
}

// ── Private helpers ───────────────────────────────────────────────────────────

func (s *Service) call(ctx context.Context, prompt string) (string, error) {
	body := anthropicRequest{
		Model:     s.Model,
		MaxTokens: 1024,
		Messages:  []anthropicMessage{{Role: "user", Content: prompt}},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("reasoning: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPI, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("reasoning: build request: %w", err)
	}
	req.Header.Set("x-api-key", s.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("reasoning: http: %w", err)
	}
	defer resp.Body.Close()

	var apiResp anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("reasoning: decode response: %w", err)
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("reasoning: api error [%s]: %s",
			apiResp.Error.Type, apiResp.Error.Message)
	}

	if len(apiResp.Content) == 0 || apiResp.Content[0].Type != "text" {
		return "", fmt.Errorf("reasoning: unexpected empty content from api")
	}

	return apiResp.Content[0].Text, nil
}

func parseOutput(raw string) (*reasoningOutput, error) {
	// Strip markdown fences if the model ignores the "no markdown" instruction.
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```") {
		lines := strings.Split(cleaned, "\n")
		cleaned = strings.Join(lines[1:len(lines)-1], "\n")
	}

	var out reasoningOutput
	if err := json.Unmarshal([]byte(cleaned), &out); err != nil {
		preview := cleaned
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		return nil, fmt.Errorf("reasoning: parse json: %w\nraw response: %s", err, preview)
	}
	return &out, nil
}

func assembleRecord(
	agentID, sessionID, intent string,
	out *reasoningOutput,
	pol models.PolicyResult,
) *models.DecisionRecord {
	status := "APPROVED"
	if !pol.BudgetOK || !pol.VendorOK || pol.PriceAnomaly {
		status = "BLOCKED"
	}

	// Hash only the stable reasoning fields, not CommittedTx or Outcome.
	hashInput := struct {
		Intent        string
		VendorChosen  models.VendorOption
		Alternatives  []models.Alternative
		ExpectedValue string
		Confidence    float64
	}{intent, out.VendorChosen, out.Alternatives, out.ExpectedValue, out.Confidence}

	reasoningHash, _ := decision.HashCanonicalJSON(hashInput)

	return &models.DecisionRecord{
		AgentID:       agentID,
		SessionID:     sessionID,
		TaskIntent:    intent,
		VendorChosen:  out.VendorChosen,
		Alternatives:  out.Alternatives,
		ExpectedValue: out.ExpectedValue,
		Confidence:    out.Confidence,
		Policy:        pol,
		Status:        status,
		ReasoningHash: reasoningHash,
		Timestamp:     time.Now(),
		// CommittedTx and Outcome are set later by the orchestrator.
	}
}

// ── Prompt ────────────────────────────────────────────────────────────────────

func buildPrompt(intent string, vendors []models.VendorOption, policy models.PolicyResult) string {
	var sb strings.Builder

	sb.WriteString("You are a purchasing decision engine for an autonomous AI agent.\n")
	sb.WriteString("Choose the best vendor and explain the reasoning.\n\n")

	sb.WriteString(fmt.Sprintf("Task: %s\n\n", intent))

	sb.WriteString("Available vendors:\n")
	for _, v := range vendors {
		sb.WriteString(fmt.Sprintf(
			"- %s | price: %.4f EURQ | accuracy: %.0f%% | free: %v | notes: %s\n",
			v.Name, v.PriceEURQ, v.AccuracyScore*100, v.IsFree, v.Notes,
		))
	}

	sb.WriteString("\nPolicy result:\n")
	sb.WriteString(fmt.Sprintf("- Budget OK: %v\n", policy.BudgetOK))
	sb.WriteString(fmt.Sprintf("- Vendor on allowlist: %v\n", policy.VendorOK))
	sb.WriteString(fmt.Sprintf("- Price anomaly: %v\n", policy.PriceAnomaly))
	if policy.Reason != "" {
		sb.WriteString(fmt.Sprintf("- Reason: %s\n", policy.Reason))
	}

	sb.WriteString(`
Return ONLY valid JSON. No markdown. No explanation. No text before or after.
Exact structure required:

{
  "VendorChosen": {
    "name": "...",
    "url": "...",
    "price_eurq": 0.001,
    "is_free": false,
    "accuracy_score": 0.91,
    "notes": "..."
  },
  "Alternatives": [
    {
      "vendor": {
        "name": "...",
        "url": "...",
        "price_eurq": 0.0,
        "is_free": true,
        "accuracy_score": 0.64,
        "notes": "..."
      },
      "reason_rejected": "One sentence. Reference the actual accuracy or price data."
    }
  ],
  "ExpectedValue": "One sentence describing the concrete benefit expected.",
  "Confidence": 0.87
}

Rules:
- VendorChosen must be copied exactly from the vendor list above (same name, url, etc.)
- List every non-chosen vendor in Alternatives with a specific, data-driven rejection reason
- Confidence is 0.0-1.0 reflecting how certain you are the chosen vendor is the right call
- If policy blocked a vendor (allowlist or anomaly), that vendor must appear in Alternatives
  with reason_rejected explaining the policy issue
`)

	return sb.String()
}
