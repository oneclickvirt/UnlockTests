package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oneclickvirt/UnlockTests/executor"
)

func TestMergeMetadataPreservesAliasesAndAddsCurrentProviders(t *testing.T) {
	providers := mergeMetadata([]string{"MetaAI", "NewAI"}, map[string]string{"metaai": "ai", "newai": "ai"}, []executor.ProviderMetadata{{Name: "Meta AI", Category: "ai", Aliases: []string{"MetaAI"}}, {Name: "Removed", Category: "ai"}})
	if len(providers) != 2 || providers[0].Name != "Meta AI" || providers[0].Aliases[0] != "MetaAI" || providers[1].Name != "NewAI" {
		t.Fatalf("unexpected merged metadata: %#v", providers)
	}
	if providers[1].Category != "ai" {
		t.Fatalf("new provider category = %q", providers[1].Category)
	}
}

func TestUpdateSnapshotDoesNotRewriteSemanticMatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metadata.json")
	original := []byte(`{"schema_version":"goecs.unlocktests/provider-metadata-v1","providers":[{"name":"B","category":"ai","aliases":["two","one"]},{"name":"A","category":"ai"}]}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := updateSnapshot(path, []string{"A", "B"}, map[string]string{"a": "ai", "b": "ai"}, 1); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("semantic match rewrote snapshot: %s", after)
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
