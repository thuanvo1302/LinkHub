package password

import (
	"crypto/sha256"
	"encoding/hex"
)

// This scaffold keeps dependencies light for local bootstrap.
// Replace with bcrypt or argon2 before production use.
func Hash(raw string) string {
	sum := sha256.Sum256([]byte("linkhub::" + raw))
	return hex.EncodeToString(sum[:])
}

func Compare(raw, hashed string) bool {
	return Hash(raw) == hashed
}
