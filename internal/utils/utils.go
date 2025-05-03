package utils

import (
	"strings"
	"unicode"

	"github.com/yassirdeveloper/cli/errors"
)

func ValidateSQLName(s string) errors.Error {
	s = strings.TrimSpace(s)

	if s == "" {
		return errors.New("cannot be empty")
	}

	if unicode.IsDigit(rune(s[0])) {
		return errors.New("cannot start with a digit")
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return errors.New("cannot include a special character")
		}
	}

	reservedKeywords := []string{"SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "IN", "IS", "NULL", "TRUE", "FALSE"}
	for _, keyword := range reservedKeywords {
		if strings.ToUpper(s) == keyword {
			return errors.New("cannot be a reserved keyword")
		}
	}

	return nil
}
