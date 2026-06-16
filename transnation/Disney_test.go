package transnation

import "testing"

func TestDisneyAuthorizations(t *testing.T) {
	deviceAuth, tokenAuth := disneyAuthorizations("token-value")
	if deviceAuth != "Bearer token-value" || tokenAuth != "token-value" {
		t.Fatalf("raw token normalized incorrectly: device=%q token=%q", deviceAuth, tokenAuth)
	}

	deviceAuth, tokenAuth = disneyAuthorizations("Bearer token-value")
	if deviceAuth != "Bearer token-value" || tokenAuth != "token-value" {
		t.Fatalf("bearer token normalized incorrectly: device=%q token=%q", deviceAuth, tokenAuth)
	}
}
