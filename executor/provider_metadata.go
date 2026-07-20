package executor

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const providerMetadataSchema = "goecs.unlocktests/provider-metadata-v1"

var providerMetadataURLs = []string{
	"https://cdn.spiritlhl.net/https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/executor/data/provider-metadata.json",
	"https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/executor/data/provider-metadata.json",
}

//go:embed data/provider-metadata.json
var embeddedProviderMetadata []byte

type ProviderMetadata struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Aliases  []string `json:"aliases,omitempty"`
}

type ProviderMetadataSource struct {
	Source   string `json:"source"`
	URL      string `json:"url,omitempty"`
	Fallback bool   `json:"fallback"`
}

type providerMetadataDocument struct {
	SchemaVersion string             `json:"schema_version"`
	Providers     []ProviderMetadata `json:"providers"`
}

// LoadProviderMetadata loads descriptive metadata only. It never changes the
// provider function registry or any private provider request logic.
func LoadProviderMetadata(ctx context.Context, client *http.Client) ([]ProviderMetadata, ProviderMetadataSource, error) {
	return loadProviderMetadata(ctx, client, providerMetadataURLs, embeddedProviderMetadata)
}

func loadProviderMetadata(ctx context.Context, client *http.Client, urls []string, embedded []byte) ([]ProviderMetadata, ProviderMetadataSource, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	embeddedProviders, err := parseProviderMetadata(embedded, 0)
	if err != nil {
		return nil, ProviderMetadataSource{}, fmt.Errorf("invalid embedded provider metadata: %w", err)
	}
	if client == nil {
		client = &http.Client{Timeout: 6 * time.Second}
	}
	for _, rawURL := range urls {
		data, fetchErr := fetchProviderMetadata(ctx, client, rawURL)
		if fetchErr != nil {
			continue
		}
		providers, parseErr := parseProviderMetadata(data, len(embeddedProviders))
		if parseErr == nil {
			return providers, ProviderMetadataSource{Source: "remote", URL: rawURL}, nil
		}
	}
	return embeddedProviders, ProviderMetadataSource{Source: "embedded", Fallback: true}, nil
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
	if document.SchemaVersion != providerMetadataSchema {
		return nil, fmt.Errorf("unsupported schema %q", document.SchemaVersion)
	}
	seen := make(map[string]struct{}, len(document.Providers))
	providers := make([]ProviderMetadata, 0, len(document.Providers))
	for _, provider := range document.Providers {
		provider.Name = strings.TrimSpace(provider.Name)
		provider.Category = strings.TrimSpace(provider.Category)
		if provider.Name == "" || provider.Category == "" {
			continue
		}
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
