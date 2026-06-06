package decision

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// HashCanonicalJSON returns the SHA-256 hex digest of JSON-marshaled input.
func HashCanonicalJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
