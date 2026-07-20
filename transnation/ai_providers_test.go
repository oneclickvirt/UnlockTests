package transnation

import (
	"net/http"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestAIProviderMetadata(t *testing.T) {
	tests := map[string]func(*http.Client) model.Result{
		"Claude":        Claude,
		"Coze":          Coze,
		"DeepSeek":      DeepSeek,
		"Dola AI":       Dola,
		"Grok":          Grok,
		"Kimi":          Kimi,
		"Mistral AI":    MistralAI,
		"Perplexity AI": PerplexityAI,
		"Poe":           Poe,
		"Sora":          Sora,
	}
	for want, fn := range tests {
		got := fn(nil)
		if got.Name != want {
			t.Fatalf("expected provider name %q, got %q", want, got.Name)
		}
		if got.Status != "" {
			t.Fatalf("expected nil-client metadata probe for %q to leave status empty, got %q", want, got.Status)
		}
	}
}

func TestXProviderMetadata(t *testing.T) {
	got := X(nil)
	if got.Name != "X (formerly Twitter)" || got.Status != "" {
		t.Fatalf("unexpected X metadata result: %#v", got)
	}
}
