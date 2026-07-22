package main

import (
	"regexp"
	"strings"
)

var privateErrorPattern = regexp.MustCompile("(?i)(https?://)[^\\s]+")
var privateTokenPattern = regexp.MustCompile("(?i)(^|[?&\\s])(token|api[_-]?key|secret|authorization|password|key)=([^&\\s]+)")
var privatePathPattern = regexp.MustCompile("(?:[A-Za-z]:\\\\|/)[^\\s]+")

func indentLegacyOutput(value string) string {
	lines := strings.Split(value, "\n")
	for index, line := range lines {
		if line != "" && line[0] != ' ' && line[0] != '\t' {
			lines[index] = " " + line
		}
	}
	return strings.Join(lines, "\n")
}

// sanitizeErrorText keeps remote-source and credential details out of user-visible errors.
func sanitizeErrorText(value string) string {
	value = privateErrorPattern.ReplaceAllString(value, "<remote>")
	value = privateTokenPattern.ReplaceAllString(value, "$1$2=<redacted>")
	return privatePathPattern.ReplaceAllString(value, "<path>")
}
