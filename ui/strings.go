package ui

import (
	"strings"
)

func StripNewLines(in string) string {
	return strings.ReplaceAll(in, "\n", " ")
}
