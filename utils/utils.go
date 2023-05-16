package utils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
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

func GetMatch(acc string, str string) bool {
	regexPattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(str))
	re := regexp.MustCompile(regexPattern)
	match := re.MatchString(acc)
	return match
}

func RewokeAccess(acc string, str string) string {
	regexPattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(str))
	re := regexp.MustCompile(regexPattern)
	replacedStr := re.ReplaceAllString(acc, "")
	return replacedStr
}
