package ui

import (
	"encoding/json"
	"io"
)

// RenderJson provides a standardized implementation for rendering JSON
func RenderJson(data interface{}, writer io.Writer) error {
	outBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	outBytes = append(outBytes, "\n"...)
	_, err = writer.Write(outBytes)
	return err
}
