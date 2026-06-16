package us

import "testing"

func TestSplitBearerAuthorization(t *testing.T) {
	bearer, raw := splitBearerAuthorization("token-value")
	if bearer != "Bearer token-value" || raw != "token-value" {
		t.Fatalf("raw token normalized incorrectly: bearer=%q raw=%q", bearer, raw)
	}

	bearer, raw = splitBearerAuthorization("Bearer token-value")
	if bearer != "Bearer token-value" || raw != "token-value" {
		t.Fatalf("bearer token normalized incorrectly: bearer=%q raw=%q", bearer, raw)
	}
}
