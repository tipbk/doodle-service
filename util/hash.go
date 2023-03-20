package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func HashPassword(password string) string {
	salt := os.Getenv("HASH_SECRET")
	passwordInput := fmt.Sprintf("%v%v", password, salt)
	hash := sha256.Sum256([]byte(passwordInput))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
