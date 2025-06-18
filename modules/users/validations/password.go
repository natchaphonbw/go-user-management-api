package validations

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
	var (
		hasUpper = false
		hasLower = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		}
	}

	if !hasUpper || !hasLower {
		return errors.New("password must contain at least one upper and lower case")
	}

	return nil

}
