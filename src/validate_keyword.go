package source

import (
	"strings"
)

func ValidateKeyword(rawKeyword string) (string, error) {
	// Trim newline characters from the input
	rawKeyword = strings.TrimSpace(rawKeyword)

	validKeyword := strings.ReplaceAll(rawKeyword, " ", "+")
	return validKeyword, nil
}
