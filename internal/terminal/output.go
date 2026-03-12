package terminal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	blue   = color.New(color.FgBlue)
	cyan   = color.New(color.FgCyan)
	bold   = color.New(color.Bold)
)

// Output writes to stderr (stdout is reserved for MCP stdio transport).
func Output(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Success prints a green checkmark message to stderr.
func Success(format string, a ...any) {
	green.Fprint(os.Stderr, "✓ ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Error prints a red X message to stderr.
func Error(format string, a ...any) {
	red.Fprint(os.Stderr, "✗ ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Warn prints a yellow warning to stderr.
func Warn(format string, a ...any) {
	yellow.Fprint(os.Stderr, "! ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Info prints a blue info message to stderr.
func Info(format string, a ...any) {
	blue.Fprint(os.Stderr, "i ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

// Bold prints bold text to stderr.
func Bold(format string, a ...any) {
	bold.Fprintf(os.Stderr, format+"\n", a...)
}

// Cyan returns a cyan-colored string.
func Cyan(s string) string {
	return cyan.Sprint(s)
}

// Red returns a red-colored string.
func Red(s string) string {
	return red.Sprint(s)
}

// Green returns a green-colored string.
func Green(s string) string {
	return green.Sprint(s)
}

// Yellow returns a yellow-colored string.
func Yellow(s string) string {
	return yellow.Sprint(s)
}

// MaskString masks a string after the first n visible characters.
func MaskString(s string, visible int) string {
	if len(s) <= visible {
		return s
	}
	masked := s[:visible]
	for i := visible; i < len(s); i++ {
		masked += "•"
	}
	return masked
}
