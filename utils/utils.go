package utils

import (
	"crypto/sha256"
	"time"
)

// Hash creates a sha256 hash of a given string
func Hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return string(bs)
}

func GetCurrentTimestamp() int {
	var date = time.Now().Unix()
	return int(date)
}
