package util

const testnetExplorerBase = "https://testnet.explorer.perawallet.app"

// TxURL returns a Pera testnet explorer link for a transaction ID.
func TxURL(txID string) string {
	return testnetExplorerBase + "/tx/" + txID
}

// AccountURL returns a Pera testnet explorer link for an address.
func AccountURL(address string) string {
	return testnetExplorerBase + "/address/" + address
}
