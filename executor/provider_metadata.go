package executor

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	ProviderMetadataSchema            = "goecs.unlocktests/provider-metadata-v1"
	ProviderMetadataManifestSchema    = "goecs.unlocktests/provider-metadata-manifest-v1"
	DefaultProviderMetadataMinimum    = 100
	DefaultProviderMetadataSyncSource = "https://raw.githubusercontent.com/HsukqiLee/MediaUnlockTest/main/pkg/providers/lists.go"
)

var providerMetadataURLs = []string{
	"https://cdn.spiritlhl.net/https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/executor/data/provider-metadata.json",
	"https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/executor/data/provider-metadata.json",
}

//go:embed data/provider-metadata.json
var embeddedProviderMetadata []byte

//go:embed data/provider-metadata.manifest.json
var embeddedProviderMetadataManifest []byte

type ProviderMetadata struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Groups       []string `json:"groups"`
	SupportsIPv6 bool     `json:"supports_ipv6"`
	Aliases      []string `json:"aliases,omitempty"`
}

type ProviderMetadataSource struct {
	Schema      string    `json:"schema"`
	Count       int       `json:"count"`
	GeneratedAt time.Time `json:"generated_at,omitempty"`
	Source      string    `json:"source"`
	Fallback    bool      `json:"fallback"`
}

type providerMetadataDocument struct {
	SchemaVersion string             `json:"schema_version"`
	GeneratedAt   time.Time          `json:"generated_at,omitempty"`
	Providers     []ProviderMetadata `json:"providers"`
}

type providerMetadataManifest struct {
	Schema      string    `json:"schema"`
	File        string    `json:"file"`
	Count       int       `json:"count"`
	SHA256      string    `json:"sha256"`
	GeneratedAt time.Time `json:"generated_at"`
}

// LoadProviderMetadata loads descriptive metadata only. It never changes the
// provider function registry or any private provider request logic.
func LoadProviderMetadata(ctx context.Context, client *http.Client) ([]ProviderMetadata, ProviderMetadataSource, error) {
	return loadProviderMetadata(ctx, client, providerMetadataURLs, embeddedProviderMetadata, embeddedProviderMetadataManifest)
}

// EmbeddedProviderMetadata returns the validated compile-time snapshot without
// performing network access.
func EmbeddedProviderMetadata() ([]ProviderMetadata, error) {
	providers, err := parseProviderMetadata(embeddedProviderMetadata, 0)
	if err != nil {
		return nil, err
	}
	if err := validateProviderMetadataManifest(embeddedProviderMetadataManifest, embeddedProviderMetadata, len(providers), providerMetadataGeneratedAt(embeddedProviderMetadata)); err != nil {
		return nil, fmt.Errorf("validate embedded provider metadata manifest: %w", err)
	}
	return providers, nil
}

// EmbeddedProviderMetadataSnapshot returns the compile-time snapshot together
// with the version metadata needed by aggregators and GUIs.
func EmbeddedProviderMetadataSnapshot() ([]ProviderMetadata, ProviderMetadataSource, error) {
	providers, err := EmbeddedProviderMetadata()
	if err != nil {
		return nil, ProviderMetadataSource{}, err
	}
	source := ProviderMetadataSource{Schema: ProviderMetadataSchema, Count: len(providers), GeneratedAt: providerMetadataGeneratedAt(embeddedProviderMetadata), Source: "embedded", Fallback: true}
	return providers, source, nil
}

func loadProviderMetadata(ctx context.Context, client *http.Client, urls []string, embedded, embeddedManifest []byte) ([]ProviderMetadata, ProviderMetadataSource, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	embeddedProviders, err := parseProviderMetadata(embedded, 0)
	if err != nil {
		return nil, ProviderMetadataSource{}, fmt.Errorf("invalid embedded provider metadata: %w", err)
	}
	embeddedGeneratedAt := providerMetadataGeneratedAt(embedded)
	if err := validateProviderMetadataManifest(embeddedManifest, embedded, len(embeddedProviders), embeddedGeneratedAt); err != nil {
		return nil, ProviderMetadataSource{}, fmt.Errorf("invalid embedded provider metadata manifest: %w", err)
	}
	if client == nil {
		client = &http.Client{Timeout: 6 * time.Second}
	}
	for index, rawURL := range urls {
		manifest, fetchErr := fetchProviderMetadata(ctx, client, providerMetadataManifestURL(rawURL))
		if fetchErr != nil {
			continue
		}
		data, fetchErr := fetchProviderMetadata(ctx, client, rawURL)
		if fetchErr != nil {
			continue
		}
		providers, parseErr := parseProviderMetadata(data, len(embeddedProviders))
		generatedAt := providerMetadataGeneratedAt(data)
		if parseErr == nil && validateProviderMetadataManifest(manifest, data, len(providers), generatedAt) == nil {
			return providers, ProviderMetadataSource{Schema: ProviderMetadataSchema, Count: len(providers), GeneratedAt: generatedAt, Source: providerRemoteSource(index), Fallback: index > 0}, nil
		}
	}
	return embeddedProviders, ProviderMetadataSource{Schema: ProviderMetadataSchema, Count: len(embeddedProviders), GeneratedAt: embeddedGeneratedAt, Source: "embedded", Fallback: true}, nil
}

func providerRemoteSource(index int) string {
	if index == 0 {
		return "cdn"
	}
	if index == 1 {
		return "raw"
	}
	return "remote"
}

func providerMetadataManifestURL(snapshotURL string) string {
	return strings.TrimSuffix(snapshotURL, ".json") + ".manifest.json"
}

func validateProviderMetadataManifest(data, snapshot []byte, count int, generatedAt time.Time) error {
	var manifest providerMetadataManifest
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("manifest contains trailing JSON")
	}
	if manifest.Schema != ProviderMetadataManifestSchema || manifest.File != "provider-metadata.json" || manifest.Count != count || manifest.GeneratedAt.IsZero() || !manifest.GeneratedAt.Equal(generatedAt) {
		return errors.New("manifest schema, file, count, or generated_at is invalid")
	}
	hash := sha256.Sum256(snapshot)
	if !strings.EqualFold(manifest.SHA256, hex.EncodeToString(hash[:])) {
		return errors.New("manifest SHA-256 does not match snapshot")
	}
	return nil
}

func providerMetadataGeneratedAt(data []byte) time.Time {
	var document struct {
		GeneratedAt time.Time `json:"generated_at"`
	}
	if json.Unmarshal(data, &document) != nil {
		return time.Time{}
	}
	return document.GeneratedAt
}

func fetchProviderMetadata(ctx context.Context, client *http.Client, rawURL string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", response.StatusCode)
	}
	const maximumSize = 4 << 20
	data, err := io.ReadAll(io.LimitReader(response.Body, maximumSize+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maximumSize {
		return nil, fmt.Errorf("provider metadata exceeds %d bytes", maximumSize)
	}
	return data, nil
}

func parseProviderMetadata(data []byte, minimum int) ([]ProviderMetadata, error) {
	var document providerMetadataDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	if document.SchemaVersion != ProviderMetadataSchema {
		return nil, fmt.Errorf("unsupported schema %q", document.SchemaVersion)
	}
	seen := make(map[string]struct{}, len(document.Providers))
	seenIDs := make(map[string]struct{}, len(document.Providers))
	providers := make([]ProviderMetadata, 0, len(document.Providers))
	for _, provider := range document.Providers {
		provider.Name = strings.TrimSpace(provider.Name)
		provider.Category = strings.TrimSpace(provider.Category)
		if provider.Name == "" || provider.Category == "" {
			continue
		}
		provider.ID = strings.TrimSpace(provider.ID)
		if provider.ID == "" {
			provider.ID = providerMetadataID(provider.Name)
		}
		if provider.ID == "" {
			continue
		}
		if _, exists := seenIDs[provider.ID]; exists {
			return nil, fmt.Errorf("duplicate provider id %q", provider.ID)
		}
		seenIDs[provider.ID] = struct{}{}
		provider.Groups = normalizeProviderGroups(provider.Groups, provider.Category)
		key := strings.ToLower(provider.Name)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		aliases := provider.Aliases[:0]
		for _, alias := range provider.Aliases {
			alias = strings.TrimSpace(alias)
			if alias != "" {
				aliases = append(aliases, alias)
			}
		}
		provider.Aliases = aliases
		providers = append(providers, provider)
	}
	if len(providers) == 0 || len(providers) < minimum {
		return nil, fmt.Errorf("provider metadata count %d is below minimum %d", len(providers), minimum)
	}
	sort.Slice(providers, func(i, j int) bool { return strings.ToLower(providers[i].Name) < strings.ToLower(providers[j].Name) })
	return providers, nil
}

func providerMetadataID(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var builder strings.Builder
	dash := false
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			builder.WriteRune(char)
			dash = false
			continue
		}
		if builder.Len() > 0 && !dash {
			builder.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func normalizeProviderGroups(groups []string, category string) []string {
	seen := make(map[string]struct{}, len(groups)+1)
	for _, group := range groups {
		group = strings.TrimSpace(strings.ToLower(group))
		if group != "" {
			seen[group] = struct{}{}
		}
	}
	if len(seen) == 0 {
		if group := providerCategoryGroup(category); group != "" {
			seen[group] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for group := range seen {
		result = append(result, group)
	}
	sort.Strings(result)
	return result
}

func providerCategoryGroup(category string) string {
	switch strings.TrimSpace(strings.ToLower(category)) {
	case "global":
		return "global"
	case "hong-kong":
		return "hk"
	case "taiwan":
		return "tw"
	case "japan":
		return "jp"
	case "korea":
		return "kr"
	case "north-america":
		return "na"
	case "south-america":
		return "sa"
	case "europe":
		return "eu"
	case "africa":
		return "africa"
	case "south-east-asia":
		return "sea"
	case "oceania":
		return "oceania"
	case "sports":
		return "sports"
	case "ai":
		return "ai"
	default:
		return ""
	}
}
