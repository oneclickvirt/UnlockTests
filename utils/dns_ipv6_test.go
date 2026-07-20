package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestNormalizeResultFillsNameAndDNSStatus(t *testing.T) {
	err := &net.DNSError{Err: "no such host", Name: "example.invalid", IsNotFound: true}
	result := NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Status: model.StatusNetworkErr, Err: err},
		"Fallback",
	)
	if result.Name != "Fallback" {
		t.Fatalf("expected fallback name, got %q", result.Name)
	}
	if result.Status != model.StatusDNSFailed {
		t.Fatalf("expected %s, got %s", model.StatusDNSFailed, result.Status)
	}
	if result.Err != err {
		t.Fatalf("expected original error to be preserved")
	}
}

func TestNormalizeResultKeepsNilNetworkErrorUnified(t *testing.T) {
	result := NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "Test", Status: model.StatusNetworkErr},
		"Fallback",
	)
	if result.Name != "Test" {
		t.Fatalf("expected original name, got %q", result.Name)
	}
	if result.Status != model.StatusNetworkErr {
		t.Fatalf("expected %s, got %s", model.StatusNetworkErr, result.Status)
	}
	if result.Err != nil {
		t.Fatalf("expected nil error to stay nil")
	}
}

func TestNormalizeResultDetectsIPv6NoAddressWithCustomClient(t *testing.T) {
	originalLookupIP := lookupIP
	lookupIP = func(ctx context.Context, network, hostname string) ([]net.IP, error) {
		switch network {
		case "ip4":
			return []net.IP{net.IPv4(203, 0, 113, 10)}, nil
		case "ip6":
			return nil, &net.DNSError{Err: "no such host", Name: hostname, IsNotFound: true}
		default:
			return nil, errors.New("unexpected network")
		}
	}
	defer func() { lookupIP = originalLookupIP }()
	SetDNSIPVersion("ipv6")
	defer SetDNSIPVersion("")

	err := &url.Error{
		Op:  "Get",
		URL: "https://ipv6.example.test/path",
		Err: errors.New("no suitable address found"),
	}
	result := NormalizeResult(
		&http.Client{Transport: &http.Transport{}},
		model.Result{Status: model.StatusNetworkErr, Err: err},
		"IPv6 Service",
	)
	if result.Status != model.StatusNoIPv6 {
		t.Fatalf("expected %s, got %s", model.StatusNoIPv6, result.Status)
	}
	if result.Name != "IPv6 Service" {
		t.Fatalf("expected fallback name, got %q", result.Name)
	}
}

func TestIsIPv6ClientRecognizesCallerContextWrapper(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := WithCallerContext(Ipv6HttpClient, ctx)
	if !IsIPv6Client(client) {
		t.Fatal("caller context wrapper hid IPv6 transport identity")
	}
}

func TestNormalizeResultFillsUnexpectedForEmptyStatus(t *testing.T) {
	result := NormalizeResult(nil, model.Result{}, "Fallback")
	if result.Name != "Fallback" {
		t.Fatalf("expected fallback name, got %q", result.Name)
	}
	if result.Status != model.StatusUnexpected {
		t.Fatalf("expected %s, got %s", model.StatusUnexpected, result.Status)
	}
}

func TestNormalizeResultMarksWAFTimeoutAsBanned(t *testing.T) {
	result := NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "Slow", Status: model.StatusNetworkErr, Err: context.DeadlineExceeded},
		"Fallback",
	)
	if result.Status != model.StatusBanned {
		t.Fatalf("expected %s, got %s", model.StatusBanned, result.Status)
	}

	result = NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "Slow", Status: model.StatusNetworkErr, Err: &timeoutErr{}},
		"Fallback",
	)
	if result.Status != model.StatusBanned {
		t.Fatalf("expected net.Error timeout to become %s, got %s", model.StatusBanned, result.Status)
	}
}

func TestNormalizeResultMarksWAFStatusErrorAsBanned(t *testing.T) {
	result := NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "Blocked", Status: model.StatusUnexpected, Err: errors.New("get blocked failed with code: 503")},
		"Fallback",
	)
	if result.Status != model.StatusBanned {
		t.Fatalf("expected %s, got %s", model.StatusBanned, result.Status)
	}
}

func TestNormalizeResultMarksUnavailableStatusErrorAsNo(t *testing.T) {
	for _, code := range []int{400, 404, 451, 452} {
		result := NormalizeResult(
			&http.Client{Transport: Ipv4Transport},
			model.Result{Name: "Unavailable", Status: model.StatusUnexpected, Err: errors.New("get service failed with code: " + strconv.Itoa(code))},
			"Fallback",
		)
		if result.Status != model.StatusNo {
			t.Fatalf("expected status code %d to become %s, got %s", code, model.StatusNo, result.Status)
		}
	}
}

func TestNormalizeResultUnifiesManualUnavailableStatuses(t *testing.T) {
	result := NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "Cloudflare", Status: model.StatusNo, Info: "Banned by cloudflare"},
		"Fallback",
	)
	if result.Status != model.StatusBanned {
		t.Fatalf("expected %s, got %s", model.StatusBanned, result.Status)
	}

	result = NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "RateLimit", Status: model.StatusNo, Info: "429 Rate limit"},
		"Fallback",
	)
	if result.Status != model.StatusRateLimited {
		t.Fatalf("expected %s, got %s", model.StatusRateLimited, result.Status)
	}

	result = NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "RateLimit", Status: model.StatusUnexpected, Info: "Too Many Requests"},
		"Fallback",
	)
	if result.Status != model.StatusRateLimited {
		t.Fatalf("expected Too Many Requests to become %s, got %s", model.StatusRateLimited, result.Status)
	}

	result = NormalizeResult(
		&http.Client{Transport: Ipv4Transport},
		model.Result{Name: "RateLimit", Status: model.StatusUnexpected, Err: errors.New("request failed with code: 429")},
		"Fallback",
	)
	if result.Status != model.StatusRateLimited {
		t.Fatalf("expected 429 error to become %s, got %s", model.StatusRateLimited, result.Status)
	}
}

func TestNormalizeResultMapsRateLimitForAnyProviderStatus(t *testing.T) {
	for _, status := range []string{model.StatusErr, model.StatusNetworkErr, model.StatusUnexpected, model.StatusNo} {
		got := NormalizeResult(nil, model.Result{Name: "fixture", Status: status, Err: fmt.Errorf("HTTP status 429: rate limited")}, "fixture")
		if got.Status != model.StatusRateLimited {
			t.Fatalf("status %q normalized to %q", status, got.Status)
		}
	}
}

type timeoutErr struct{}

func (e *timeoutErr) Error() string   { return "request timed out" }
func (e *timeoutErr) Timeout() bool   { return true }
func (e *timeoutErr) Temporary() bool { return true }

func TestSetCustomDNSServersNormalizesHostPort(t *testing.T) {
	SetCustomDNSServers("1.1.1.1:53 [2606:4700:4700::1111]:53")
	defer SetCustomDNSServers("")
	got := get_nameserver_from_resolv()
	if len(got) != 2 {
		t.Fatalf("expected 2 custom DNS servers, got %d", len(got))
	}
	if got[0] != "1.1.1.1" || got[1] != "2606:4700:4700::1111" {
		t.Fatalf("unexpected normalized DNS servers: %#v", got)
	}
}
