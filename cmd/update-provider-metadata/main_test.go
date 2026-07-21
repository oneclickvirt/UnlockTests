package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/oneclickvirt/UnlockTests/executor"
)

func TestParseReferenceProvidersUsesGoSyntaxAndDeduplicates(t *testing.T) {
	providers, err := parseReferenceProviders([]byte(`package providers
var GlobeTests = []TestItem{{"Alpha", Alpha, true}, {"Shared", Shared, false}}
var JapanTests = []TestItem{{"Beta", Beta, false}, {"Shared", SharedJP, false}}
var Helper = []string{"not a provider"}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(providers) != 3 || providers[0].Name != "Alpha" || providers[0].Category != "global" || len(providers[0].Groups) != 1 || providers[0].Groups[0] != "global" || !providers[0].SupportsIPv6 || providers[1].Name != "Beta" {
		t.Fatalf("unexpected providers: %#v", providers)
	}
}

func TestFetchReferenceProvidersRejectsOversizedOrBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusBadGateway)
	}))
	defer server.Close()
	if _, err := fetchReferenceProviders(context.Background(), server.Client(), server.URL); err == nil {
		t.Fatal("bad status unexpectedly accepted")
	}
}

func TestMergeMetadataPreservesAliasesAndAddsCurrentProviders(t *testing.T) {
	reference := map[string]referenceProvider{"metaai": {Name: "MetaAI", Category: "ai", Groups: []string{"ai"}, SupportsIPv6: true}}
	providers := mergeMetadata([]string{"MetaAI", "NewAI"}, map[string]string{"metaai": "ai", "newai": "ai"}, reference, []executor.ProviderMetadata{{Name: "Meta AI", Category: "ai", Aliases: []string{"MetaAI"}}, {Name: "Removed", Category: "ai"}})
	if len(providers) != 2 || providers[0].Name != "Meta AI" || providers[0].Aliases[0] != "MetaAI" || providers[1].Name != "NewAI" {
		t.Fatalf("unexpected merged metadata: %#v", providers)
	}
	if providers[1].Category != "ai" {
		t.Fatalf("new provider category = %q", providers[1].Category)
	}
	if providers[0].ID != "metaai" || len(providers[0].Groups) != 1 || providers[0].Groups[0] != "ai" || !providers[0].SupportsIPv6 {
		t.Fatalf("metadata contract not populated: %#v", providers[0])
	}
}

func TestMergeMetadataUsesReferenceIDForCompatibleAlias(t *testing.T) {
	reference := map[string]referenceProvider{
		"microsoft copilot": {Name: "Microsoft Copilot", Category: "ai", Groups: []string{"ai"}, SupportsIPv6: true},
	}
	current := []executor.ProviderMetadata{{Name: "Copilot", Category: "ai", Aliases: []string{"Microsoft Copilot"}}}
	providers := mergeMetadata([]string{"Microsoft Copilot"}, map[string]string{"microsoft copilot": "ai"}, reference, current)
	if len(providers) != 1 || providers[0].ID != "microsoft-copilot" || !providers[0].SupportsIPv6 {
		t.Fatalf("reference-compatible ID was not retained: %#v", providers)
	}
}

func TestEnsureUniqueProviderIDsResolvesSlugCollisions(t *testing.T) {
	providers := ensureUniqueProviderIDs([]executor.ProviderMetadata{
		{ID: "fod-fuji-tv", Name: "FOD(Fuji TV)"},
		{ID: "fod-fuji-tv", Name: "FOD (Fuji TV)"},
	})
	if providers[0].ID == providers[1].ID {
		t.Fatalf("duplicate IDs were retained: %#v", providers)
	}
}

func TestUpdateSnapshotMigratesLegacyMetadataContract(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metadata.json")
	original := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","generated_at":"2026-07-20T00:00:00Z","providers":[{"name":"B","category":"ai","aliases":["two","one"]},{"name":"A","category":"ai"}]}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := updateSnapshot(path, []string{"A", "B"}, map[string]string{"a": "ai", "b": "ai"}, nil, 1); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var migrated dataDocument
	if err := json.Unmarshal(after, &migrated); err != nil {
		t.Fatal(err)
	}
	if migrated.GeneratedAt.IsZero() || len(migrated.Providers) != 2 || migrated.Providers[0].ID == "" || len(migrated.Providers[0].Groups) == 0 {
		t.Fatalf("legacy metadata contract was not populated: %s", after)
	}
	manifestData, err := os.ReadFile(filepath.Join(filepath.Dir(path), "metadata.manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest dataManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil || manifest.Schema != executor.ProviderMetadataManifestSchema || manifest.Count != 2 || manifest.SHA256 == "" {
		t.Fatalf("manifest was not generated: %s", manifestData)
	}
}

func TestProviderCatalogIncludesAllAndPreservesAICategory(t *testing.T) {
	names, categories, err := providerCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if len(names) < 100 || categories["dola ai"] != "ai" {
		t.Fatalf("unexpected provider catalog: count=%d dola=%q", len(names), categories["dola ai"])
	}
}
