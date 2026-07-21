package executor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadProviderMetadataPrefersValidatedRemote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metadata.json" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"name":"Zeta","category":"ai"},{"name":"Alpha","category":"ai"}]}`))
	}))
	defer server.Close()
	embedded := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","providers":[{"name":"Embedded","category":"ai"}]}`)
	providers, source, err := loadProviderMetadata(context.Background(), server.Client(), []string{server.URL + "/metadata.json"}, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if source.Source != "remote" || source.Fallback || source.Schema != ProviderMetadataSchema || source.Count != 2 || source.GeneratedAt.IsZero() || len(providers) != 2 || providers[0].Name != "Alpha" {
		t.Fatalf("unexpected metadata result: %#v %#v", providers, source)
	}
}

func TestLoadProviderMetadataFallsBackOnSchemaFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"schema_version":"wrong","providers":[{"name":"Remote","category":"ai"}]}`))
	}))
	defer server.Close()
	embedded := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","providers":[{"name":"Embedded","category":"ai"}]}`)
	providers, source, err := loadProviderMetadata(context.Background(), server.Client(), []string{server.URL}, embedded)
	if err != nil {
		t.Fatal(err)
	}
	if source.Source != "embedded" || !source.Fallback || len(providers) != 1 || providers[0].Name != "Embedded" {
		t.Fatalf("unexpected fallback metadata: %#v %#v", providers, source)
	}
}

func TestEmbeddedProviderMetadataCoversAIRegistry(t *testing.T) {
	providers, err := EmbeddedProviderMetadata()
	if err != nil {
		t.Fatal(err)
	}
	known := make(map[string]struct{})
	for _, provider := range providers {
		known[strings.ToLower(provider.Name)] = struct{}{}
		for _, alias := range provider.Aliases {
			known[strings.ToLower(alias)] = struct{}{}
		}
	}
	names, err := ListPlatforms("21")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if _, ok := known[strings.ToLower(name)]; !ok {
			t.Errorf("AI provider %q has no embedded metadata", name)
		}
	}
}
