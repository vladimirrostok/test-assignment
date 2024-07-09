package word

import (
	"fmt"
	domain_errors "test-assignment/domain/errors"
	"unicode"
)

// Verify the word, only a full string of unicode letters is allowed.
func IsValidUnicode(s string) error {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("error: %w", domain_errors.InvalidWordData{})
		}
	}
	return nil
}

// Verify the word, only a certain length is allowed.
func IsValidLength(s string, length int) error {
	if len([]rune(s)) != length {
		return fmt.Errorf("error: %w", domain_errors.InvalidWordLength{})
	}

	return nil
}
