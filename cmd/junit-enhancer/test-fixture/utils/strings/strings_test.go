package strings

import (
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"12345", "54321"},
		{"Hello World", "dlroW olleH"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := Reverse(test.input)
			if result != test.expected {
				t.Errorf("Reverse(%q) = %q, expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestIsPalindrome(t *testing.T) {
	palindromes := []string{
		"racecar",
		"A man a plan a canal Panama",
		"race a car", // This is NOT a palindrome but tests the cleaning
		"",
		"a",
	}

	notPalindromes := []string{
		"hello",
		"world",
		"test123",
	}

	t.Run("palindromes", func(t *testing.T) {
		for _, s := range palindromes {
			// Note: "race a car" is actually NOT a palindrome, so we'll test both cases
			if s == "race a car" {
				continue // Skip this one as it's not actually a palindrome
			}
			if !IsPalindrome(s) {
				t.Errorf("IsPalindrome(%q) should be true", s)
			}
		}
	})

	t.Run("not_palindromes", func(t *testing.T) {
		notPalindromes = append(notPalindromes, "race a car") // Add the non-palindrome
		for _, s := range notPalindromes {
			if IsPalindrome(s) {
				t.Errorf("IsPalindrome(%q) should be false", s)
			}
		}
	})
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  hello   world  ", 2},
		{"one two three four five", 5},
		{"\t\n  \t\n", 0},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := CountWords(test.input)
			if result != test.expected {
				t.Errorf("CountWords(%q) = %d, expected %d", test.input, result, test.expected)
			}
		})
	}
}

func FuzzReverse(f *testing.F) {
	f.Add("hello")
	f.Add("")
	f.Add("a")
	f.Add("hello world")
	
	f.Fuzz(func(t *testing.T, s string) {
		reversed := Reverse(s)
		doubleReversed := Reverse(reversed)
		if s != doubleReversed {
			t.Errorf("Double reverse of %q should equal original, got %q", s, doubleReversed)
		}
	})
}
