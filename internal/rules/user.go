package rules

import (
	"regexp"
	"unicode"
)

// IsValidUserName checks whether a name is valid. Expects it to be already trimmed
//
// For a name to be valid, it must have between 3 and 32 characters and only lower/uppercase letters and digits
func IsValidUserName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
	return re.MatchString(name)
}

// IsValidUserPassword checks whether a password is valid. For a password to be valid, it can't be longer than 72
// characters and must have at least:
//
//   - 8 characters
//   - an uppercase letter
//   - a lowercase letter
//   - a digit
//   - a special symbol, such as !@#$%¨&*()
func IsValidUserPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}

	var upp, low, num, sym bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsNumber(char):
			num = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
		default:
			return false
		}
	}

	return upp && low && num && sym
}
