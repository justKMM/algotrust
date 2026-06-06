package models

// AlgorandSpikeResult captures output from the Algorand commitment spike.
type AlgorandSpikeResult struct {
	Address       string
	MicroAlgos    uint64
	ReasoningHash string
	TxID          string
	ExplorerURL   string
}

// X402ProbeResult captures the outcome of an unpaid x402 probe request.
type X402ProbeResult struct {
	URL             string
	StatusCode      int
	PaymentRequired bool
	PaymentHeader   string
	BodySnippet     string
}
