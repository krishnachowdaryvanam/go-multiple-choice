package validator

import (
	"strings"
)

// isValidOption checks if the provided answer option is valid.
func IsValidOption(answer string) bool {
	// You can customize this function based on your input validation rules.
	validOptions := []string{"A", "B", "C", "D"}
	for _, option := range validOptions {
		if strings.EqualFold(answer, option) {
			return true
		}
	}
	return false
}
