package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(strLength int) string {
	if strLength <= 0 {
		return ""
	}

	result := make([]byte, strLength)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < strLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			result[i] = charset[i%len(charset)]
		} else {
			result[i] = charset[randomIndex.Int64()]
		}
	}

	return string(result)
}
