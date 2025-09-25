package color

import "fmt"

// ANSI color codes
const (
	Reset   = "\033[0m"
	Blue    = "\033[34m"
	Cyan    = "\033[36m"
	Green   = "\033[32m"
	Magenta = "\033[35m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
)

// Colorize wraps text with a color code
func Colorize(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}
