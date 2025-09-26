package utils

import "unicode"

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var hasUpper, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasNumber && hasSpecial
}
