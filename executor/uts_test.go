package executor

import (
	"net/http"
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

func TestParseSelectionHandlesRepeatedSpaces(t *testing.T) {
	if !parseSelection("0   10") {
		t.Fatalf("expected selection to parse")
	}
	if !M || !TW {
		t.Fatalf("expected global and Taiwan selections to be enabled")
	}
}

func TestParseSelectionRejectsInvalidWithoutKeepingOldState(t *testing.T) {
	if !parseSelection("0") {
		t.Fatalf("expected initial selection to parse")
	}
	if parseSelection("0 invalid") {
		t.Fatalf("expected invalid selection to be rejected")
	}
	if M || TW || HK || JP || KR || NA || SA || EU || AFR || OCEA || SPORT {
		t.Fatalf("expected invalid selection to reset all selection flags")
	}
}

func TestUniqueFuncListDeduplicatesByResultName(t *testing.T) {
	a := func(c *http.Client) model.Result { return model.Result{Name: "A"} }
	b := func(c *http.Client) model.Result { return model.Result{Name: "B"} }
	dupeA := func(c *http.Client) model.Result { return model.Result{Name: "A"} }
	got := uniqueFuncList([](func(c *http.Client) model.Result){a, b, dupeA})
	if len(got) != 2 {
		t.Fatalf("expected 2 unique funcs, got %d", len(got))
	}
}

func TestFinallyPrintResultIPv6UsesSelectedPlatformTitle(t *testing.T) {
	resetOptions()
	defer func() {
		resetOptions()
		Names = nil
		R = nil
	}()
	TW = true
	Names = []string{"Example"}
	R = []*model.Result{{Name: "Example", Status: model.StatusYes}}
	got := finallyPrintResult("en", "ipv6")
	if !strings.Contains(got, "[ Taiwan ]") {
		t.Fatalf("expected IPv6 output to use selected platform title, got %q", got)
	}
}
