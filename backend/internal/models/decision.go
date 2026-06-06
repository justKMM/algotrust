package models

import "time"

// DecisionStatus is the policy outcome for a spend request.
type DecisionStatus string

const (
	StatusApproved DecisionStatus = "APPROVED"
	StatusBlocked  DecisionStatus = "BLOCKED"
	StatusPending  DecisionStatus = "PENDING"
)

// LegacyAlternative describes a vendor the agent considered but did not choose (used by the API/frontend model).
type LegacyAlternative struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// VendorOption describes a candidate vendor for an agent task.
type VendorOption struct {
	Name          string  `json:"name"`
	URL           string  `json:"url"`
	PriceEURQ     float64 `json:"price_eurq"`
	IsFree        bool    `json:"is_free"`
	AccuracyScore float64 `json:"accuracy_score"`
	Notes         string  `json:"notes"`
}

// PolicyResult is the output of the policy engine for a spend request.
type PolicyResult struct {
	BudgetOK     bool    `json:"budget_ok"`
	VendorOK     bool    `json:"vendor_ok"`
	PriceAnomaly bool    `json:"price_anomaly"`
	AnomalyRatio float64 `json:"anomaly_ratio"`
	Reason       string  `json:"reason"`
}

// Alternative describes a vendor the agent considered but rejected (used in the service layer model).
type Alternative struct {
	Vendor         VendorOption `json:"vendor"`
	ReasonRejected string       `json:"reason_rejected"`
}

// OutcomeRecord records post-purchase results vs expectations.
type OutcomeRecord struct {
	Predicted   string  `json:"predicted"`
	Actual      string  `json:"actual"`
	Score       float64 `json:"score"`
	GroundTruth string  `json:"ground_truth"`
}

// DecisionRecord is the full audit record for an agent spend request (service layer model).
type DecisionRecord struct {
	ID            string         `json:"id"`
	AgentID       string         `json:"agent_id"`
	SessionID     string         `json:"session_id"`
	TaskIntent    string         `json:"task_intent"`
	VendorChosen  VendorOption   `json:"vendor_chosen"`
	Alternatives  []Alternative  `json:"alternatives"`
	ExpectedValue string         `json:"expected_value"`
	Confidence    float64        `json:"confidence"`
	Policy        PolicyResult   `json:"policy"`
	Status        string         `json:"status"`
	ReasoningHash string         `json:"reasoning_hash"`
	CommittedTx   string         `json:"committed_tx,omitempty"`
	Outcome       *OutcomeRecord `json:"outcome,omitempty"`
	Timestamp     time.Time      `json:"timestamp"`
}

// PolicyChecks captures policy engine results for a decision.
type PolicyChecks struct {
	BudgetOk      bool    `json:"budgetOk"`
	Reputation    float64 `json:"reputation"`
	Anomaly       string  `json:"anomaly"`
	VendorAllowed bool    `json:"vendorAllowed"`
}

// Outcome records post-purchase results vs expectations.
type Outcome struct {
	Predicted  string  `json:"predicted"`
	Actual     string  `json:"actual"`
	Verdict    string  `json:"verdict"`
	TrustDelta float64 `json:"trustDelta"`
}

// Decision is the audit record for an agent spend request (API/frontend model).
type Decision struct {
	ID            string              `json:"id"`
	Vendor        string              `json:"vendor"`
	Status        DecisionStatus      `json:"status"`
	AmountEURQ    float64             `json:"amountEURQ"`
	Intent        string              `json:"intent"`
	Alternatives  []LegacyAlternative `json:"alternatives"`
	ExpectedValue string              `json:"expectedValue"`
	Confidence    float64             `json:"confidence"`
	Policy        PolicyChecks        `json:"policy"`
	ReasoningHash string              `json:"reasoningHash"`
	Round         int64               `json:"round"`
	Timestamp     int64               `json:"timestamp"`
	Outcome       *Outcome            `json:"outcome,omitempty"`
	BlockedReason string              `json:"blockedReason,omitempty"`
	CommittedTx   string              `json:"committedTx,omitempty"`
	ExplorerURL   string              `json:"explorerUrl,omitempty"`
}

// Vendor tracks trust score for a vendor.
type Vendor struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

// Alert is a policy or anomaly notification.
type Alert struct {
	ID      string `json:"id"`
	Level   string `json:"level"`
	Message string `json:"message"`
	At      int64  `json:"at"`
}

// AppState is the dashboard snapshot served to the frontend.
type AppState struct {
	Agent          string     `json:"agent"`
	Balance        float64    `json:"balance"`
	Spent          float64    `json:"spent"`
	DailyLimit     float64    `json:"dailyLimit"`
	Decisions      []Decision `json:"decisions"`
	Vendors        []Vendor   `json:"vendors"`
	AllowedVendors []string   `json:"allowedVendors"`
	BlockedVendors []string   `json:"blockedVendors"`
	Alerts         []Alert    `json:"alerts"`
	SelectedID     *string    `json:"selectedId"`
}
