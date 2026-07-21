package executor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoadProviderMetadataPrefersValidatedRemote(t *testing.T) {
	remote := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"name":"Zeta","category":"ai"},{"name":"Alpha","category":"ai"}]}`)
	manifest := providerTestManifest(remote, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".manifest.json") {
			_, _ = w.Write(manifest)
			return
		}
		if r.URL.Path != "/metadata.json" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(remote)
	}))
	defer server.Close()
	embedded := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"name":"Embedded","category":"ai"}]}`)
	providers, source, err := loadProviderMetadata(context.Background(), server.Client(), []string{server.URL + "/metadata.json"}, embedded, providerTestManifest(embedded, 1))
	if err != nil {
		t.Fatal(err)
	}
	if source.Source != "cdn" || source.Fallback || source.Schema != ProviderMetadataSchema || source.Count != 2 || source.GeneratedAt.IsZero() || len(providers) != 2 || providers[0].Name != "Alpha" {
		t.Fatalf("unexpected metadata result: %#v %#v", providers, source)
	}
}

func TestLoadProviderMetadataFallsBackOnSchemaFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"schema_version":"wrong","providers":[{"name":"Remote","category":"ai"}]}`))
	}))
	defer server.Close()
	embedded := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"name":"Embedded","category":"ai"}]}`)
	providers, source, err := loadProviderMetadata(context.Background(), server.Client(), []string{server.URL + "/metadata.json"}, embedded, providerTestManifest(embedded, 1))
	if err != nil {
		t.Fatal(err)
	}
	if source.Source != "embedded" || !source.Fallback || len(providers) != 1 || providers[0].Name != "Embedded" {
		t.Fatalf("unexpected fallback metadata: %#v %#v", providers, source)
	}
}

func TestLoadProviderMetadataFallsBackFromCDNToRaw(t *testing.T) {
	payload := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"id":"one","name":"One","category":"global","groups":["global"],"supports_ipv6":true}]}`)
	validManifest := providerTestManifest(payload, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/cdn/metadata.manifest.json" {
			_, _ = writer.Write([]byte(`{"schema":"goecs.unlocktests/provider-metadata-manifest-v1","file":"provider-metadata.json","count":1,"sha256":"bad","generated_at":"2026-07-20T00:00:00Z"}`))
			return
		}
		if request.URL.Path == "/raw/metadata.manifest.json" {
			_, _ = writer.Write(validManifest)
			return
		}
		_, _ = writer.Write(payload)
	}))
	defer server.Close()
	providers, source, err := loadProviderMetadata(context.Background(), server.Client(), []string{server.URL + "/cdn/metadata.json", server.URL + "/raw/metadata.json"}, payload, validManifest)
	if err != nil || len(providers) != 1 || source.Source != "raw" || !source.Fallback {
		t.Fatalf("unexpected raw fallback: providers=%#v source=%#v err=%v", providers, source, err)
	}
}

func providerTestManifest(snapshot []byte, count int) []byte {
	hash := sha256.Sum256(snapshot)
	return []byte(fmt.Sprintf(`{"schema":"goecs.unlocktests/provider-metadata-manifest-v1","file":"provider-metadata.json","count":%d,"sha256":"%s","generated_at":"2026-07-20T00:00:00Z"}`, count, hex.EncodeToString(hash[:])))
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

func TestEmbeddedProviderMetadataHasCompleteLegacyContract(t *testing.T) {
	providers, err := EmbeddedProviderMetadata()
	if err != nil {
		t.Fatal(err)
	}
	ids := make(map[string]struct{}, len(providers))
	for _, provider := range providers {
		if provider.ID == "" || len(provider.Groups) == 0 {
			t.Fatalf("provider contract is incomplete: %#v", provider)
		}
		if _, exists := ids[provider.ID]; exists {
			t.Fatalf("duplicate provider id %q", provider.ID)
		}
		ids[provider.ID] = struct{}{}
	}
}
