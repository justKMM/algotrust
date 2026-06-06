package main

import (
	"fmt"
	"os"

	"rationalgo/internal/config"
	algosvc "rationalgo/internal/services/algorand"
	x402svc "rationalgo/internal/services/x402"
	"rationalgo/internal/util"
)

func main() {
	config.LoadEnv()

	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}

	if len(os.Args) < 2 {
		printStatus(cfg)
		return
	}

	switch os.Args[1] {
	case "status":
		printStatus(cfg)
	case "spike":
		runSpike(cfg, os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func runSpike(cfg config.Config, args []string) {
	target := "all"
	if len(args) > 0 {
		target = args[0]
	}

	switch target {
	case "algorand":
		spikeAlgorand(cfg)
	case "x402":
		spikeX402(cfg)
	case "all":
		spikeAlgorand(cfg)
		fmt.Println()
		spikeX402(cfg)
	default:
		fail(fmt.Errorf("unknown spike target %q (use algorand, x402, or all)", target))
	}
}

func spikeAlgorand(cfg config.Config) {
	fmt.Println("=== Algorand spike (hash commitment) ===")
	svc, err := algosvc.NewService(cfg)
	if err != nil {
		fail(err)
	}
	result, err := svc.RunSpike()
	if err != nil {
		fail(err)
	}
	fmt.Printf("wallet:         %s\n", result.Address)
	fmt.Printf("balance:        %d microAlgos\n", result.MicroAlgos)
	fmt.Printf("reasoning_hash: %s\n", result.ReasoningHash)
	fmt.Printf("tx_id:          %s\n", result.TxID)
	fmt.Printf("explorer:       %s\n", result.ExplorerURL)
	fmt.Println("ok: testnet commitment confirmed")
}

func spikeX402(cfg config.Config) {
	fmt.Println("=== x402 spike (402 probe) ===")
	result, err := x402svc.NewService(cfg).RunProbe()
	if err != nil {
		fail(err)
	}
	fmt.Printf("url:         %s\n", result.URL)
	fmt.Printf("status:      %d\n", result.StatusCode)
	if result.PaymentRequired {
		fmt.Println("payment:     HTTP 402 Payment Required (expected)")
	} else {
		fmt.Println("payment:     no 402 — endpoint may have changed; check URL")
	}
	if result.PaymentHeader != "" {
		fmt.Printf("header:      PAYMENT-REQUIRED present (%d chars)\n", len(result.PaymentHeader))
	}
	if result.BodySnippet != "" {
		fmt.Printf("body:        %s\n", result.BodySnippet)
	}
	fmt.Println("ok: x402 probe complete (paid flow lands in Phase 2)")
}

func printStatus(cfg config.Config) {
	fmt.Println("RationAlgo — Phase 0")
	fmt.Println()
	fmt.Printf("wallet:      %s\n", displayWallet(cfg))
	fmt.Printf("algod:       %s\n", cfg.AlgodURL)
	fmt.Printf("algod token: %s\n", displayToken(cfg.AlgodToken))
	fmt.Printf("x402 probe:  %s\n", cfg.X402ProbeURL)
	fmt.Printf("http addr:   %s (Phase 1+)\n", cfg.HTTPAddr)
	fmt.Println()
	if err := cfg.ValidateForSpike(); err != nil {
		fmt.Printf("spike ready: no — %v\n", err)
		fmt.Println()
		printUsage()
		return
	}
	fmt.Printf("account:     %s\n", util.AccountURL(cfg.WalletAddress))
	fmt.Println("spike ready: yes")
	fmt.Println()
	printUsage()
}

func displayWallet(cfg config.Config) string {
	if !cfg.WalletConfigured() {
		return "(not set — paste RATIONALGO_WALLET_ADDRESS in backend/.env)"
	}
	return cfg.WalletAddress
}

func displayToken(token string) string {
	if token == "" {
		return "(empty — OK for public AlgoNode)"
	}
	return "(set)"
}

func printUsage() {
	fmt.Println(`Usage:
  go run ./cmd/rationalgo                 # config status
  go run ./cmd/rationalgo spike all       # Algorand commit + x402 probe
  go run ./cmd/rationalgo spike algorand  # testnet hash commitment only
  go run ./cmd/rationalgo spike x402      # unpaid 402 probe only

Setup:
  cp .env.example .env
  # edit backend/.env — wallet address, algod token, mnemonic`)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
