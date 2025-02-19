package tool

import (
	"fmt"
)

// Colorize outputs using [ANSI Escape Codes](https://en.wikipedia.org/wiki/ANSI_escape_code)

func color(code int, in string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", code, in)
}

func FormatBold(in string) string {
	return color(1, in)
}

func FormatRed(in string) string {
	return color(31, in)
}

func FormatGreen(in string) string {
	return color(32, in)
}

func FormatYellow(in string) string {
	return color(33, in)
}

func FormatBrightBlack(in string) string {
	return color(90, in)
}
