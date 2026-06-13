package slug

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

var reserved = map[string]struct{}{
	"admin":     {},
	"api":       {},
	"login":     {},
	"register":  {},
	"dashboard": {},
}

func Normalize(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func Reserved(value string) bool {
	_, exists := reserved[Normalize(value)]
	return exists
}

func Random(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(alphabet[r.Intn(len(alphabet))])
	}
	return b.String()
}
