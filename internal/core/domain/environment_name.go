package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrEnvironmentNameEmpty = errors.New("environment name must not be empty")

// NormalizeEnvironmentName trims the supplied environment name and validates it
// against the documented character rules. The resulting string preserves case
// and internal whitespace while rejecting path separators and other unsupported
// characters.
func NormalizeEnvironmentName(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ErrEnvironmentNameEmpty
	}

	for _, r := range trimmed {
		switch {
		case r == '/' || r == '\\':
			return "", fmt.Errorf("environment name %q must not contain path separators", raw)
		case unicode.IsLetter(r):
			continue
		case unicode.IsDigit(r):
			continue
		case unicode.IsSpace(r):
			continue
		case r == '-' || r == '_' || r == '.':
			continue
		default:
			return "", fmt.Errorf("environment name %q contains unsupported character %q", raw, r)
		}
	}

	return trimmed, nil
}
