package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateRandomHash(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
