package transnation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestDolaParsesRegion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<script>window.__data={"inner_region":"US"}</script>`))
	}))
	defer server.Close()

	got := checkDola(server.Client(), server.URL+"/chat/", "127.0.0.1")
	if got.Name != "Dola AI" || got.Status != model.StatusYes || got.Region != "us" {
		t.Fatalf("unexpected Dola result: %#v", got)
	}
}

func TestDolaStatusMapping(t *testing.T) {
	tests := map[string]struct {
		code int
		want string
	}{
		"forbidden":    {code: http.StatusForbidden, want: model.StatusBanned},
		"rate limited": {code: http.StatusTooManyRequests, want: model.StatusRateLimited},
		"malformed":    {code: http.StatusOK, want: model.StatusUnexpected},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.code)
				_, _ = w.Write([]byte(`{"missing":"region"}`))
			}))
			defer server.Close()

			got := checkDola(server.Client(), server.URL, "127.0.0.1")
			if got.Status != tt.want {
				t.Fatalf("got %#v, want status %q", got, tt.want)
			}
		})
	}
}

func TestXStatusMapping(t *testing.T) {
	tests := map[string]struct {
		region string
		code   int
		body   string
		want   string
	}{
		"available":          {region: "US", code: http.StatusOK, want: model.StatusYes},
		"restricted country": {region: "CN", code: http.StatusOK, want: model.StatusNo},
		"forbidden":          {region: "US", code: http.StatusForbidden, want: model.StatusNo},
		"rate limited":       {region: "US", code: http.StatusTooManyRequests, want: model.StatusRateLimited},
		"waf":                {region: "US", code: http.StatusOK, body: "Attention Required", want: model.StatusBanned},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/cdn-cgi/trace" {
					_, _ = w.Write([]byte("loc=" + tt.region + "\n"))
					return
				}
				w.WriteHeader(tt.code)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			got := checkX(server.Client(), server.URL+"/", server.URL+"/cdn-cgi/trace", "127.0.0.1")
			if got.Status != tt.want {
				t.Fatalf("got %#v, want status %q", got, tt.want)
			}
			if got.Region != strings.ToLower(tt.region) {
				t.Fatalf("got region %q, want %q", got.Region, strings.ToLower(tt.region))
			}
		})
	}
}

func TestXTraceRateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cdn-cgi/trace" {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		t.Fatal("main X endpoint must not run after trace rate limit")
	}))
	defer server.Close()

	got := checkX(server.Client(), server.URL+"/", server.URL+"/cdn-cgi/trace", "127.0.0.1")
	if got.Status != model.StatusRateLimited || got.Info != "trace HTTP 429" {
		t.Fatalf("unexpected result: %#v", got)
	}
}
