package policy

import (
	"fmt"

	"rationalgo/internal/models"
)

// Evaluate runs three sequential policy checks against the proposed spend.
// It stops at the first failure and returns the PolicyResult.
func Evaluate(
	chosen models.VendorOption,
	amountEURQ float64,
	dailySpent float64,
	dailyLimit float64,
	allowedVendors []string,
	priceHistory map[string][]float64,
) models.PolicyResult {
	// 1. Budget check
	if dailySpent+amountEURQ > dailyLimit {
		return models.PolicyResult{
			BudgetOK: false,
			VendorOK: true,
			Reason:   fmt.Sprintf("daily limit exceeded: spent %.4f of %.4f EURQ", dailySpent, dailyLimit),
		}
	}

	// 2. Vendor allowlist check
	allowed := false
	for _, v := range allowedVendors {
		if v == chosen.Name {
			allowed = true
			break
		}
	}
	if !allowed {
		return models.PolicyResult{
			BudgetOK: true,
			VendorOK: false,
			Reason:   "vendor not in allowlist",
		}
	}

	// 3. Price anomaly check
	history := priceHistory[chosen.Name]
	if len(history) > 0 {
		var sum float64
		for _, p := range history {
			sum += p
		}
		avg := sum / float64(len(history))
		if avg > 0 && amountEURQ > 3*avg {
			ratio := amountEURQ / avg
			return models.PolicyResult{
				BudgetOK:     true,
				VendorOK:     true,
				PriceAnomaly: true,
				AnomalyRatio: ratio,
				Reason:       fmt.Sprintf("price %.1fx above 7-day average", ratio),
			}
		}
	}

	return models.PolicyResult{
		BudgetOK: true,
		VendorOK: true,
	}
}
