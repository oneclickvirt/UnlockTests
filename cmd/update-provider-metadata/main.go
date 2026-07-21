package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/executor"
)

const schemaVersion = executor.ProviderMetadataSchema

type dataDocument struct {
	SchemaVersion string                      `json:"schema_version"`
	GeneratedAt   time.Time                   `json:"generated_at"`
	Providers     []executor.ProviderMetadata `json:"providers"`
}

type dataManifest struct {
	Schema      string    `json:"schema"`
	File        string    `json:"file"`
	Count       int       `json:"count"`
	SHA256      string    `json:"sha256"`
	GeneratedAt time.Time `json:"generated_at"`
}

func main() {
	output := flag.String("output", "executor/data/provider-metadata.json", "snapshot output path")
	source := flag.String("source", executor.DefaultProviderMetadataSyncSource, "reference provider registry URL")
	timeout := flag.Duration("timeout", 30*time.Second, "reference provider fetch timeout")
	minimum := flag.Int("min-count", executor.DefaultProviderMetadataMinimum, "minimum accepted provider count")
	minimumSource := flag.Int("min-source-count", executor.DefaultProviderMetadataMinimum, "minimum accepted reference provider count")
	flag.Parse()
	names, categories, err := providerCatalog()
	if err != nil {
		fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	reference, err := fetchReferenceProviders(ctx, http.DefaultClient, *source)
	if err != nil {
		fatal(err)
	}
	if len(reference) < *minimumSource {
		fatal(fmt.Errorf("reference provider count %d is below minimum %d", len(reference), *minimumSource))
	}
	referenceByName := make(map[string]referenceProvider, len(reference))
	for _, provider := range reference {
		key := strings.ToLower(provider.Name)
		referenceByName[key] = provider
		if _, exists := categories[key]; !exists {
			categories[key] = provider.Category
			names = append(names, provider.Name)
		}
	}
	if err := updateSnapshot(*output, names, categories, referenceByName, *minimum); err != nil {
		fatal(err)
	}
}

type referenceProvider struct {
	Name         string
	Category     string
	Groups       []string
	SupportsIPv6 bool
}

func fetchReferenceProviders(ctx context.Context, client *http.Client, source string) ([]referenceProvider, error) {
	if client == nil {
		client = http.DefaultClient
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "oneclickvirt-unlocktests-provider-sync/1")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reference provider registry returned HTTP %d", response.StatusCode)
	}
	const maximumSize = 4 << 20
	data, err := io.ReadAll(io.LimitReader(response.Body, maximumSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maximumSize {
		return nil, fmt.Errorf("reference provider registry exceeds %d bytes", maximumSize)
	}
	return parseReferenceProviders(data)
}

func parseReferenceProviders(data []byte) ([]referenceProvider, error) {
	file, err := parser.ParseFile(token.NewFileSet(), "lists.go", data, 0)
	if err != nil {
		return nil, fmt.Errorf("parse reference provider registry: %w", err)
	}
	unique := make(map[string]referenceProvider)
	for _, declaration := range file.Decls {
		general, ok := declaration.(*ast.GenDecl)
		if !ok || general.Tok != token.VAR {
			continue
		}
		for _, rawSpec := range general.Specs {
			spec, ok := rawSpec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for index, name := range spec.Names {
				if index >= len(spec.Values) {
					continue
				}
				category := referenceCategory(name.Name)
				group := referenceGroup(name.Name)
				if category == "" {
					continue
				}
				list, ok := spec.Values[index].(*ast.CompositeLit)
				if !ok {
					continue
				}
				for _, rawItem := range list.Elts {
					item, ok := rawItem.(*ast.CompositeLit)
					if !ok || len(item.Elts) < 2 || isNilProviderExpression(item.Elts[1]) {
						continue
					}
					literal, ok := item.Elts[0].(*ast.BasicLit)
					if !ok || literal.Kind != token.STRING {
						continue
					}
					providerName, unquoteErr := strconv.Unquote(literal.Value)
					providerName = strings.TrimSpace(providerName)
					if unquoteErr != nil || providerName == "" {
						continue
					}
					key := strings.ToLower(providerName)
					provider := unique[key]
					if provider.Name == "" {
						provider.Name = providerName
						provider.Category = category
					}
					provider.Groups = appendUniqueString(provider.Groups, group)
					if len(item.Elts) >= 3 {
						if enabled, ok := item.Elts[2].(*ast.Ident); ok && enabled.Name == "true" {
							provider.SupportsIPv6 = true
						}
					}
					unique[key] = provider
				}
			}
		}
	}
	providers := make([]referenceProvider, 0, len(unique))
	for _, provider := range unique {
		sort.Strings(provider.Groups)
		providers = append(providers, provider)
	}
	sort.Slice(providers, func(i, j int) bool { return strings.ToLower(providers[i].Name) < strings.ToLower(providers[j].Name) })
	if len(providers) == 0 {
		return nil, fmt.Errorf("reference provider registry contains no providers")
	}
	return providers, nil
}

func isNilProviderExpression(expression ast.Expr) bool {
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == "nil"
}

func appendUniqueString(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, current := range values {
		if current == value {
			return values
		}
	}
	return append(values, value)
}

func referenceGroup(name string) string {
	groups := map[string]string{
		"GlobeTests": "global", "HongKongTests": "hk", "TaiwanTests": "tw",
		"JapanTests": "jp", "KoreaTests": "kr", "NorthAmericaTests": "na",
		"SouthAmericaTests": "sa", "EuropeTests": "eu", "AfricaTests": "africa",
		"SouthEastAsiaTests": "sea", "OceaniaTests": "oceania", "SportsTests": "sports", "AITests": "ai",
	}
	return groups[name]
}

func referenceCategory(name string) string {
	name = strings.TrimSuffix(name, "Tests")
	switch name {
	case "Globe":
		return "global"
	case "HongKong":
		return "hong-kong"
	case "Taiwan":
		return "taiwan"
	case "Japan":
		return "japan"
	case "Korea":
		return "korea"
	case "NorthAmerica":
		return "north-america"
	case "SouthAmerica":
		return "south-america"
	case "SouthEastAsia":
		return "south-east-asia"
	case "Europe":
		return "europe"
	case "Africa":
		return "africa"
	case "Oceania":
		return "oceania"
	case "Sports":
		return "sports"
	case "AI":
		return "ai"
	default:
		return ""
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

func updateSnapshot(path string, names []string, categories map[string]string, reference map[string]referenceProvider, minimum int) error {
	currentData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var current dataDocument
	if err := json.Unmarshal(currentData, &current); err != nil || current.SchemaVersion != schemaVersion {
		return fmt.Errorf("current provider metadata schema is invalid")
	}
	nextProviders := mergeMetadata(names, categories, reference, current.Providers)
	if len(nextProviders) < minimum {
		return fmt.Errorf("provider count %d is below minimum %d", len(nextProviders), minimum)
	}
	if len(current.Providers) > 0 && len(nextProviders)*10 < len(current.Providers)*7 {
		return fmt.Errorf("provider count dropped from %d to %d", len(current.Providers), len(nextProviders))
	}
	if sameMetadata(current.Providers, nextProviders) && !current.GeneratedAt.IsZero() {
		return writeProviderManifest(path, currentData, len(current.Providers))
	}
	next, err := json.MarshalIndent(dataDocument{SchemaVersion: schemaVersion, GeneratedAt: time.Now().UTC(), Providers: nextProviders}, "", "  ")
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
	if err := os.Rename(temporaryName, path); err != nil {
		return err
	}
	return writeProviderManifest(path, next, len(nextProviders))
}

func writeProviderManifest(snapshotPath string, snapshot []byte, count int) error {
	var document dataDocument
	if err := json.Unmarshal(snapshot, &document); err != nil {
		return err
	}
	hash := sha256.Sum256(snapshot)
	manifest := dataManifest{
		Schema: executor.ProviderMetadataManifestSchema, File: filepath.Base(snapshotPath), Count: count,
		SHA256: hex.EncodeToString(hash[:]), GeneratedAt: document.GeneratedAt,
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	manifestPath := strings.TrimSuffix(snapshotPath, ".json") + ".manifest.json"
	if current, readErr := os.ReadFile(manifestPath); readErr == nil && string(current) == string(data) {
		return nil
	}
	return atomicWriteFile(manifestPath, data)
}

func atomicWriteFile(path string, data []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".provider-manifest-*.json")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(data); err != nil {
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
	return os.Rename(name, path)
}

func sameMetadata(left, right []executor.ProviderMetadata) bool {
	normalize := func(values []executor.ProviderMetadata) []executor.ProviderMetadata {
		result := make([]executor.ProviderMetadata, 0, len(values))
		for _, provider := range values {
			provider.Name = strings.TrimSpace(provider.Name)
			provider.Category = strings.TrimSpace(provider.Category)
			provider.ID = strings.TrimSpace(provider.ID)
			provider.Groups = append([]string(nil), provider.Groups...)
			for index := range provider.Groups {
				provider.Groups[index] = strings.TrimSpace(provider.Groups[index])
			}
			sort.Strings(provider.Groups)
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
		if a[index].ID != b[index].ID || a[index].Name != b[index].Name || a[index].Category != b[index].Category || a[index].SupportsIPv6 != b[index].SupportsIPv6 || len(a[index].Groups) != len(b[index].Groups) || len(a[index].Aliases) != len(b[index].Aliases) {
			return false
		}
		for groupIndex := range a[index].Groups {
			if a[index].Groups[groupIndex] != b[index].Groups[groupIndex] {
				return false
			}
		}
		for aliasIndex := range a[index].Aliases {
			if a[index].Aliases[aliasIndex] != b[index].Aliases[aliasIndex] {
				return false
			}
		}
	}
	return true
}

func mergeMetadata(names []string, categories map[string]string, reference map[string]referenceProvider, current []executor.ProviderMetadata) []executor.ProviderMetadata {
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
			provider = enrichProviderMetadata(provider, categories[strings.ToLower(name)], reference[strings.ToLower(name)])
			result = append(result, provider)
		} else {
			category := categories[strings.ToLower(name)]
			if category == "" {
				category = "other"
			}
			provider := executor.ProviderMetadata{Name: name, Category: category}
			result = append(result, enrichProviderMetadata(provider, category, reference[strings.ToLower(name)]))
		}
	}
	return ensureUniqueProviderIDs(result)
}

func ensureUniqueProviderIDs(providers []executor.ProviderMetadata) []executor.ProviderMetadata {
	// Prefer the most descriptive spelling for the base ID, then assign stable
	// numeric suffixes to aliases whose slug would otherwise collide.
	sort.SliceStable(providers, func(i, j int) bool {
		if providers[i].ID != providers[j].ID {
			return providers[i].ID < providers[j].ID
		}
		if len(providers[i].Name) != len(providers[j].Name) {
			return len(providers[i].Name) > len(providers[j].Name)
		}
		return providers[i].Name < providers[j].Name
	})
	used := make(map[string]struct{}, len(providers))
	for index := range providers {
		base := providers[index].ID
		if base == "" {
			base = metadataID(providers[index].Name)
		}
		candidate := base
		for suffix := 2; ; suffix++ {
			if _, exists := used[candidate]; !exists {
				break
			}
			candidate = fmt.Sprintf("%s-%d", base, suffix)
		}
		providers[index].ID = candidate
		used[candidate] = struct{}{}
	}
	sort.Slice(providers, func(i, j int) bool { return strings.ToLower(providers[i].Name) < strings.ToLower(providers[j].Name) })
	return providers
}

func enrichProviderMetadata(provider executor.ProviderMetadata, category string, reference referenceProvider) executor.ProviderMetadata {
	if strings.TrimSpace(provider.Category) == "" {
		provider.Category = category
	}
	currentID := strings.TrimSpace(provider.ID)
	if currentID == "" || reference.Name != "" && currentID == metadataID(provider.Name) {
		idName := provider.Name
		if reference.Name != "" {
			idName = reference.Name
		}
		provider.ID = metadataID(idName)
	}
	groups := make(map[string]struct{}, len(provider.Groups)+len(reference.Groups)+1)
	for _, group := range append(append([]string(nil), provider.Groups...), reference.Groups...) {
		group = strings.TrimSpace(strings.ToLower(group))
		if group != "" {
			groups[group] = struct{}{}
		}
	}
	if len(groups) == 0 {
		if group := categoryGroup(provider.Category); group != "" {
			groups[group] = struct{}{}
		}
	}
	provider.Groups = provider.Groups[:0]
	for group := range groups {
		provider.Groups = append(provider.Groups, group)
	}
	sort.Strings(provider.Groups)
	provider.SupportsIPv6 = provider.SupportsIPv6 || reference.SupportsIPv6
	return provider
}

func metadataID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			builder.WriteRune(char)
			lastDash = false
			continue
		}
		if !lastDash && builder.Len() > 0 {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func categoryGroup(category string) string {
	groups := map[string]string{
		"global": "global", "hong-kong": "hk", "taiwan": "tw", "japan": "jp", "korea": "kr",
		"north-america": "na", "south-america": "sa", "europe": "eu", "africa": "africa",
		"south-east-asia": "sea", "oceania": "oceania", "sports": "sports", "ai": "ai",
	}
	return groups[strings.ToLower(strings.TrimSpace(category))]
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
