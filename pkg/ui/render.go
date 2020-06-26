package ui

import (
	"io"
	"os"
)

const (
	// OutputFormatHuman prints the output for humans
	OutputFormatHuman = "human"
	// OutputFormatJson prints the output in JSON
	OutputFormatJson = "json"
)

var outputFormat string

// RenderView renders a view to stdout. If an error occurs it will exit with a non-zero exit code.
// Use RenderViewWriter if you wish control of the io.Writer or error handling
// Use SetOutputFormat to set the global output format
func RenderView(view View) {
	ExitIfError(RenderViewWriter(view, os.Stdout))
}

func RenderViewWriter(view View, writer io.Writer) error {
	switch outputFormat {
	case OutputFormatJson:
		return view.RenderJson(writer)
	default:
		return view.RenderHuman(writer)
	}
}

// SetOutputFormat sets the global output format for all calls to RenderView* funcs
func SetOutputFormat(newOutputFormat string) {
	outputFormat = newOutputFormat
}
