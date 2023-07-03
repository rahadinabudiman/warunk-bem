package helpers

import (
	"errors"
	"unicode"
)

func ValidateUsername(username string) error {
	// check username for only alphaNumeric characters
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("only alphanumeric characters allowed for username")
		}
	}
	// check username length
	if 5 <= len(username) && len(username) <= 14 {
		return nil
	}
	return errors.New("username length must be greater than 4 and less than 15 characters")
}
