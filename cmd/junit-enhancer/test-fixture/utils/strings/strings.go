package strings

import (
	"strings"
	"unicode"
)

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	cleaned := ""
	for _, r := range s {
		if unicode.IsLetter(r) {
			cleaned += string(r)
		}
	}
	return cleaned == Reverse(cleaned)
}

func CountWords(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Fields(s))
}
