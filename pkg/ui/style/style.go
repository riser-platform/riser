package style

import (
	"fmt"

	"github.com/wzshiming/ctc"
)

func Good(message string) string {
	return Colorize(message, ctc.ForegroundBrightGreen)
}

func Emphasis(message string) string {
	return Colorize(message, ctc.ForegroundBrightCyan)
}

func Muted(message string) string {
	return Colorize(message, ctc.ForegroundBrightBlack)
}

func Bad(message string) string {
	return Colorize(message, ctc.ForegroundBrightRed)
}

func Warn(message string) string {
	return Colorize(message, ctc.ForegroundBrightYellow)
}

func Colorize(message string, color ctc.Color) string {
	return fmt.Sprintf("%s%s%s", color, message, ctc.Reset)
}
