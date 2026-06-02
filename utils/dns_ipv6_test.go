package utils

import (
	"net"
	"net/http"
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
