package util

import "github.com/sethvargo/go-password/password"

type Password struct {
	// Maximum length of the password
	Length int

	// The number of digit characters
	NumDigits int

	// The number of symbol characters to be generated
	NumSymbols int

	// Should only allow lowercase
	NoUpper bool

	// Should repeat  characters
	AllowRepeat bool
}

func (p Password) GetPassword() string {
	return password.MustGenerate(p.Length, p.NumDigits, p.NumSymbols, p.NoUpper, p.AllowRepeat)
}

var DefaultPassword = Password{
	AllowRepeat: false,
	NoUpper: false,
	NumSymbols: 1,
	NumDigits: 1,
	Length: 10,
}