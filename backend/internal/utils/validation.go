package utils

import (
	"fmt"
	"strings"
)

// ValidateRequired validates that required fields are not empty
func ValidateRequired(fieldName string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateID validates that ID is not empty
func ValidateID(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// ValidatePositiveInt validates that an integer is positive
func ValidatePositiveInt(fieldName string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0", fieldName)
	}
	return nil
}

// CombineErrors combines multiple errors into a single error
func CombineErrors(errs ...error) error {
	var validErrs []error
	for _, err := range errs {
		if err != nil {
			validErrs = append(validErrs, err)
		}
	}
	if len(validErrs) == 0 {
		return nil
	}
	if len(validErrs) == 1 {
		return validErrs[0]
	}
	messages := make([]string, len(validErrs))
	for i, err := range validErrs {
		messages[i] = err.Error()
	}
	return fmt.Errorf(strings.Join(messages, "; "))
}
