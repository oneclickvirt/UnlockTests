package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oneclickvirt/UnlockTests/executor"
)

const schemaVersion = "goecs.unlocktests/provider-metadata-v1"

type dataDocument struct {
	SchemaVersion string                      `json:"schema_version"`
	Providers     []executor.ProviderMetadata `json:"providers"`
}

func main() {
	output := flag.String("output", "executor/data/provider-metadata.json", "snapshot output path")
	minimum := flag.Int("min-count", 1, "minimum accepted provider count")
	flag.Parse()
	names, categories, err := providerCatalog()
	if err != nil {
		fatal(err)
	}
	if err := updateSnapshot(*output, names, categories, *minimum); err != nil {
		fatal(err)
	}
}

func providerCatalog() ([]string, map[string]string, error) {
	categories := make(map[string]string)
	for _, group := range []struct{ selection, category string }{
		{"0", "global"}, {"10", "taiwan"}, {"11", "hong-kong"}, {"12", "japan"},
		{"13", "korea"}, {"14", "north-america"}, {"15", "south-america"},
		{"16", "europe"}, {"17", "africa"}, {"18", "oceania"}, {"19", "sports"}, {"21", "ai"},
	} {
		names, err := executor.ListPlatforms(group.selection)
		if err != nil {
			return nil, nil, err
		}
		for _, name := range names {
			key := strings.ToLower(strings.TrimSpace(name))
			if key != "" {
				categories[key] = group.category
			}
		}
	}
	names, err := executor.ListPlatforms("20")
	return names, categories, err
}

func updateSnapshot(path string, names []string, categories map[string]string, minimum int) error {
	currentData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var current dataDocument
	if err := json.Unmarshal(currentData, &current); err != nil || current.SchemaVersion != schemaVersion {
		return fmt.Errorf("current provider metadata schema is invalid")
	}
	nextProviders := mergeMetadata(names, categories, current.Providers)
	if len(nextProviders) < minimum {
		return fmt.Errorf("provider count %d is below minimum %d", len(nextProviders), minimum)
	}
	if len(current.Providers) > 0 && len(nextProviders)*10 < len(current.Providers)*7 {
		return fmt.Errorf("provider count dropped from %d to %d", len(current.Providers), len(nextProviders))
	}
	if sameMetadata(current.Providers, nextProviders) {
		return nil
	}
	next, err := json.MarshalIndent(dataDocument{SchemaVersion: schemaVersion, Providers: nextProviders}, "", "  ")
	if err != nil {
		return err
	}
	next = append(next, '\n')
	if string(next) == string(currentData) {
		return nil
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".provider-metadata-*.json")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if _, err := temporary.Write(next); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryName, path)
}

func sameMetadata(left, right []executor.ProviderMetadata) bool {
	normalize := func(values []executor.ProviderMetadata) []executor.ProviderMetadata {
		result := make([]executor.ProviderMetadata, 0, len(values))
		for _, provider := range values {
			provider.Name = strings.TrimSpace(provider.Name)
			provider.Category = strings.TrimSpace(provider.Category)
			provider.Aliases = append([]string(nil), provider.Aliases...)
			for index := range provider.Aliases {
				provider.Aliases[index] = strings.TrimSpace(provider.Aliases[index])
			}
			sort.Strings(provider.Aliases)
			result = append(result, provider)
		}
		sort.Slice(result, func(i, j int) bool { return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name) })
		return result
	}
	a, b := normalize(left), normalize(right)
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index].Name != b[index].Name || a[index].Category != b[index].Category || len(a[index].Aliases) != len(b[index].Aliases) {
			return false
		}
		for aliasIndex := range a[index].Aliases {
			if a[index].Aliases[aliasIndex] != b[index].Aliases[aliasIndex] {
				return false
			}
		}
	}
	return true
}

func mergeMetadata(names []string, categories map[string]string, current []executor.ProviderMetadata) []executor.ProviderMetadata {
	result := make([]executor.ProviderMetadata, 0, len(names))
	used := make(map[int]struct{}, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		match := -1
		for index, provider := range current {
			if _, ok := used[index]; ok {
				continue
			}
			if strings.EqualFold(provider.Name, name) || containsAlias(provider.Aliases, name) {
				match = index
				break
			}
		}
		if match >= 0 {
			provider := current[match]
			used[match] = struct{}{}
			result = append(result, provider)
		} else {
			category := categories[strings.ToLower(name)]
			if category == "" {
				category = "other"
			}
			result = append(result, executor.ProviderMetadata{Name: name, Category: category})
		}
	}
	sort.Slice(result, func(i, j int) bool { return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name) })
	return result
}

func containsAlias(aliases []string, name string) bool {
	for _, alias := range aliases {
		if strings.EqualFold(alias, name) {
			return true
		}
	}
	return false
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
