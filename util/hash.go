package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256生成哈希值
func GetSHA256HashCode(message []byte) string {
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}
