package strutil

import "strings"

// SafeLogString marks a string that has already been sanitized and is
// safe to interpolate into a log line. Using a distinct type (rather
// than string) gives static analyzers a real signal that this value
// has crossed a sanitization boundary.
type SafeLogString string

// SanitizeForLog strips characters that could be used to forge or
// inject fake log lines (CRLF injection) when embedding a
// user/env-controlled string into a log message.
func SanitizeForLog(s string) SafeLogString {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return SafeLogString(s)
}

// Truncate shortens s to at most n runes, appending "..." if it was
// cut. Useful for bounding untrusted strings before logging or
// displaying them.
func Truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
