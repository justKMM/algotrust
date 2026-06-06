package vendor

import "rationalgo/internal/models"

// Catalog is the hard-coded vendor list for the demo.
// Accuracy figures are grounded in WMO short-range vs 24h forecast benchmarks.
var Catalog = []models.VendorOption{
	{
		Name:          "GoPlausible WeatherAPI",
		URL:           "https://example.x402.goplausible.xyz/avm/weather",
		PriceEURQ:     0.001,
		IsFree:        false,
		AccuracyScore: 0.91,
		Notes:         "Real-time, 1h precision forecast",
	},
	{
		Name:          "OpenMeteo Free",
		URL:           "https://api.open-meteo.com/v1/forecast?latitude=50.11&longitude=8.68&hourly=precipitation_probability&forecast_days=1",
		PriceEURQ:     0,
		IsFree:        true,
		AccuracyScore: 0.64,
		Notes:         "24h window, no payment required",
	},
}

// PriceHistory holds 7 recent prices per vendor for anomaly detection.
var PriceHistory = map[string][]float64{
	"GoPlausible WeatherAPI": {0.001, 0.001, 0.001, 0.0011, 0.001, 0.001, 0.001},
	"OpenMeteo Free":         {0, 0, 0, 0, 0, 0, 0},
}

// GetAll returns all vendors in the catalog.
func GetAll() []models.VendorOption {
	result := make([]models.VendorOption, len(Catalog))
	copy(result, Catalog)
	return result
}

// GetPriceHistory returns a copy of the price history map.
func GetPriceHistory() map[string][]float64 {
	result := make(map[string][]float64, len(PriceHistory))
	for k, v := range PriceHistory {
		hist := make([]float64, len(v))
		copy(hist, v)
		result[k] = hist
	}
	return result
}
