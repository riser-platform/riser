package sdk

import (
	"fmt"
	"strings"
)

// ClientError provides the error message, status code.
type ClientError struct {
	StatusCode       int
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors"`
}

func (e *ClientError) Error() string {
	if len(e.ValidationErrors) == 0 {
		return fmt.Sprintf("Error: %s", e.Message)
	} else {
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("%s:", e.Message))
		builder.WriteString("\n")
		for fieldName, errorMessage := range e.ValidationErrors {
			builder.WriteString(fmt.Sprintf(" • %s: %s\n", fieldName, errorMessage))
		}
		return builder.String()
	}
}
