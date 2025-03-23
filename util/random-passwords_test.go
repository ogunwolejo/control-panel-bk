package util

import (
	"testing"
	"unicode"
)

func TestGetPassword(t *testing.T) {
	tests := []struct {
		name       string
		password  Password
		expectsErr bool
	}{
		{
			name: "Default settings",
			password: Password{
				Length:      10,
				NumDigits:   2,
				NumSymbols:  2,
				NoUpper:     false,
				AllowRepeat: true,
			},
		},
		{
			name: "Only lowercase letters",
			password: Password{
				Length:      12,
				NumDigits:   3,
				NumSymbols:  2,
				NoUpper:     true,
				AllowRepeat: true,
			},
		},
		{
			name: "No repeating characters",
			password: Password{
				Length:      8,
				NumDigits:   2,
				NumSymbols:  1,
				NoUpper:     false,
				AllowRepeat: false,
			},
		},
		{
			name: "Only digits and symbols",
			password: Password{
				Length:      6,
				NumDigits:   6,
				NumSymbols:  0,
				NoUpper:     true,
				AllowRepeat: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generatedPassword := tc.password.GetPassword()

			if len(generatedPassword) != tc.password.Length {
				t.Errorf("Expected password length %d, got %d", tc.password.Length, len(generatedPassword))
			}

			digitCount := 0
			symbolCount := 0
			upperCount := 0
			charSet := make(map[rune]bool)

			for _, char := range generatedPassword {
				charSet[char] = true
				if unicode.IsDigit(char) {
					digitCount++
				} else if unicode.IsSymbol(char) || unicode.IsPunct(char) {
					symbolCount++
				} else if unicode.IsUpper(char) {
					upperCount++
				}
			}

			if digitCount < tc.password.NumDigits {
				t.Errorf("Expected at least %d digits, got %d", tc.password.NumDigits, digitCount)
			}

			if symbolCount < tc.password.NumSymbols {
				t.Errorf("Expected at least %d symbols, got %d", tc.password.NumSymbols, symbolCount)
			}

			if tc.password.NoUpper && upperCount > 0 {
				t.Errorf("Expected no uppercase letters, but found %d", upperCount)
			}

			if !tc.password.AllowRepeat && len(charSet) < len(generatedPassword) {
				t.Errorf("Expected no repeated characters, but found duplicates")
			}
		})
	}
}
