package executor

import (
	"strings"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestShowResultNetworkErrorWithNilErr(t *testing.T) {
	got := ShowResult(&model.Result{Name: "Test", Status: model.StatusNetworkErr})
	if !strings.Contains(got, "Network Error") {
		t.Fatalf("expected unified network error message, got %q", got)
	}
}

func TestShowResultDNSFailed(t *testing.T) {
	got := ShowResult(&model.Result{Name: "Test", Status: model.StatusDNSFailed})
	if !strings.Contains(got, "DNS Resolve Failed") {
		t.Fatalf("expected unified DNS error message, got %q", got)
	}
}
