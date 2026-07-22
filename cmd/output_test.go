package main

import (
	"strings"
	"testing"
)

func TestIndentLegacyOutputLeavesOneLeadingCell(t *testing.T) {
	got := indentLegacyOutput("alpha\n beta\n\tgamma\n")
	want := " alpha\n beta\n\tgamma\n"
	if got != want {
		t.Fatalf("indentLegacyOutput() = %q, want %q", got, want)
	}
}

func TestSanitizeErrorTextRedactsRemoteSecretsAndPaths(t *testing.T) {
	input := "load https://private.example/repo/data?token=secret from /private/cache failed: api_key=hidden"
	got := sanitizeErrorText(input)
	for _, forbidden := range []string{"private.example", "secret", "hidden", "/private/cache"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("sanitizeErrorText() leaked %q in %q", forbidden, got)
		}
	}
}
