package sdk

import (
	"fmt"
	"net/http"
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
		// TODO: Do we need the status code here?
		return fmt.Sprintf("Received HTTP %d (%s) %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
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
