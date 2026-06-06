package models

// DecisionStatus is the policy outcome for a spend request.
type DecisionStatus string

const (
	StatusApproved DecisionStatus = "APPROVED"
	StatusBlocked  DecisionStatus = "BLOCKED"
	StatusPending  DecisionStatus = "PENDING"
)

// Alternative describes a vendor the agent considered but did not choose.
type Alternative struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// PolicyChecks captures policy engine results for a decision.
type PolicyChecks struct {
	BudgetOk       bool   `json:"budgetOk"`
	Reputation     float64 `json:"reputation"`
	Anomaly        string `json:"anomaly"`
	VendorAllowed  bool   `json:"vendorAllowed"`
}

// Outcome records post-purchase results vs expectations.
type Outcome struct {
	Predicted  string  `json:"predicted"`
	Actual     string  `json:"actual"`
	Verdict    string  `json:"verdict"`
	TrustDelta float64 `json:"trustDelta"`
}

// Decision is the audit record for an agent spend request.
type Decision struct {
	ID             string         `json:"id"`
	Vendor         string         `json:"vendor"`
	Status         DecisionStatus `json:"status"`
	AmountEURQ     float64        `json:"amountEURQ"`
	Intent         string         `json:"intent"`
	Alternatives   []Alternative  `json:"alternatives"`
	ExpectedValue  string         `json:"expectedValue"`
	Confidence     float64        `json:"confidence"`
	Policy         PolicyChecks   `json:"policy"`
	ReasoningHash  string         `json:"reasoningHash"`
	Round          int64          `json:"round"`
	Timestamp      int64          `json:"timestamp"`
	Outcome        *Outcome       `json:"outcome,omitempty"`
	BlockedReason  string         `json:"blockedReason,omitempty"`
	CommittedTx    string         `json:"committedTx,omitempty"`
	ExplorerURL    string         `json:"explorerUrl,omitempty"`
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
