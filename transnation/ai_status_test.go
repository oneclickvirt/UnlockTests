package transnation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestCheckAIRegionalStatusAcceptsDeepSeek202(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cdn-cgi/trace":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("loc=US\n"))
		default:
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer server.Close()

	got := checkAIRegionalStatus(server.Client(), aiRegionalProbe{
		name:        "DeepSeek",
		hostname:    "chat.deepseek.com",
		url:         server.URL + "/",
		traceURL:    server.URL + "/cdn-cgi/trace",
		okCodes:     map[int]bool{http.StatusAccepted: true},
		bannedCodes: map[int]bool{http.StatusForbidden: true},
		wafKeywords: defaultAIWAFKeywords(),
	})
	if got.Status != model.StatusYes {
		t.Fatalf("expected 202 Accepted to resolve as Yes, got %#v", got)
	}
	if got.Region != "us" {
		t.Fatalf("expected Cloudflare trace region us, got %q", got.Region)
	}
}

func TestCheckAIRegionalStatusUsesSupportCountryList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cdn-cgi/trace":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("loc=CN\n"))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	got := checkAIRegionalStatus(server.Client(), aiRegionalProbe{
		name:             "Perplexity AI",
		hostname:         "www.perplexity.ai",
		url:              server.URL + "/",
		traceURL:         server.URL + "/cdn-cgi/trace",
		okCodes:          map[int]bool{http.StatusOK: true},
		supportCountries: []string{"us"},
		wafKeywords:      defaultAIWAFKeywords(),
	})
	if got.Status != model.StatusNo {
		t.Fatalf("expected unsupported trace region to resolve as No, got %#v", got)
	}
	if got.Region != "cn" {
		t.Fatalf("expected region cn, got %q", got.Region)
	}
}

func TestCheckAIRegionalStatusUsesRestrictedCountryList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cdn-cgi/trace":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("loc=CN\n"))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	got := checkAIRegionalStatus(server.Client(), aiRegionalProbe{
		name:                "Grok",
		hostname:            "grok.com",
		url:                 server.URL + "/",
		traceURL:            server.URL + "/cdn-cgi/trace",
		okCodes:             map[int]bool{http.StatusOK: true},
		restrictedCountries: []string{"cn"},
		wafKeywords:         defaultAIWAFKeywords(),
	})
	if got.Status != model.StatusNo {
		t.Fatalf("expected restricted trace region to resolve as No, got %#v", got)
	}
	if got.Region != "cn" {
		t.Fatalf("expected region cn, got %q", got.Region)
	}
}

func TestCheckAIRegionalStatusForbiddenCodeUsesTraceRegion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cdn-cgi/trace":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("loc=US\n"))
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer server.Close()

	got := checkAIRegionalStatus(server.Client(), aiRegionalProbe{
		name:                "Perplexity AI",
		hostname:            "www.perplexity.ai",
		url:                 server.URL + "/",
		traceURL:            server.URL + "/cdn-cgi/trace",
		forbiddenCodes:      map[int]bool{http.StatusForbidden: true},
		restrictedCountries: []string{"cn"},
		wafKeywords:         defaultAIWAFKeywords(),
	})
	if got.Status != model.StatusBanned {
		t.Fatalf("expected allowed trace region plus forbidden status to resolve as Banned, got %#v", got)
	}
	if got.Region != "us" {
		t.Fatalf("expected region us, got %q", got.Region)
	}
}
