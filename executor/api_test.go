package executor

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestListPlatformsRejectsInvalidSelection(t *testing.T) {
	if _, err := ListPlatforms("invalid"); err == nil {
		t.Fatalf("expected invalid selection error")
	}
}

func TestNormalizeIPVersion(t *testing.T) {
	tests := map[string]string{
		"":          "ipv4",
		"4":         "ipv4",
		"ipv4":      "ipv4",
		"6":         "ipv6",
		"IPV6":      "ipv6",
		"0":         "auto",
		"dualstack": "auto",
	}
	for input, want := range tests {
		got, err := normalizeIPVersion(input)
		if err != nil {
			t.Fatalf("normalizeIPVersion(%q) returned error: %v", input, err)
		}
		if got != want {
			t.Fatalf("normalizeIPVersion(%q) = %q, want %q", input, got, want)
		}
	}
	if _, err := normalizeIPVersion("ipv5"); err == nil {
		t.Fatalf("expected invalid IP version to fail")
	}
}

func TestClientWithContextDeadlineTightensTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	source := &http.Client{Timeout: 30 * time.Second, Transport: &http.Transport{}}
	got := clientWithContextDeadline(source, ctx)
	if got == source {
		t.Fatalf("expected client clone when context deadline is tighter")
	}
	if got.Transport != source.Transport {
		t.Fatalf("expected transport to be preserved")
	}
	if got.Timeout <= 0 || got.Timeout > 100*time.Millisecond {
		t.Fatalf("expected timeout to be tightened to context deadline, got %s", got.Timeout)
	}
}

func TestClientWithContextDeadlineKeepsShorterClientTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	source := &http.Client{Timeout: 10 * time.Millisecond, Transport: &http.Transport{}}
	got := clientWithContextDeadline(source, ctx)
	if got != source {
		t.Fatalf("expected existing shorter timeout client to be reused")
	}
}

func TestRunFunctionsStructuredReturnsOrderedResults(t *testing.T) {
	funcs := []func(c *http.Client) model.Result{
		func(c *http.Client) model.Result {
			if c == nil {
				return model.Result{Name: "A"}
			}
			return model.Result{Name: "A", Status: model.StatusYes, Region: "us"}
		},
		func(c *http.Client) model.Result {
			if c == nil {
				return model.Result{Name: "B"}
			}
			return model.Result{Name: "B", Status: model.StatusNo}
		},
	}
	results, err := runFunctionsStructured(context.Background(), funcs, RunOptions{Client: &http.Client{}, IPVersion: "ipv4", Concurrency: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Name != "A" || results[0].Status != model.StatusYes || results[0].Region != "us" {
		t.Fatalf("unexpected first result: %#v", results[0])
	}
	if results[1].Name != "B" || results[1].Status != model.StatusNo {
		t.Fatalf("unexpected second result: %#v", results[1])
	}
}

func TestRunFunctionsStructuredHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	funcs := []func(c *http.Client) model.Result{
		func(c *http.Client) model.Result {
			if c == nil {
				return model.Result{Name: "Canceled"}
			}
			t.Fatalf("test function should not run with a canceled context")
			return model.Result{}
		},
	}
	results, err := runFunctionsStructured(ctx, funcs, RunOptions{Client: &http.Client{}, IPVersion: "ipv4"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(results) != 1 || results[0].Name != "Canceled" || results[0].Status != model.StatusErr {
		t.Fatalf("unexpected canceled result: %#v", results)
	}
}

func TestRunFunctionWithContextConvertsPanic(t *testing.T) {
	result, err := runFunctionWithContext(context.Background(), func(c *http.Client) model.Result {
		if c == nil {
			return model.Result{Name: "Panic"}
		}
		panic("boom")
	}, RunOptions{Client: &http.Client{}, IPVersion: "ipv4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Panic" || result.Status != model.StatusErr || result.Err == nil {
		t.Fatalf("expected panic to become structured error, got %#v", result)
	}
}

func TestRunFunctionWithContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	result, err := runFunctionWithContext(ctx, func(c *http.Client) model.Result {
		if c == nil {
			return model.Result{Name: "Slow"}
		}
		time.Sleep(100 * time.Millisecond)
		return model.Result{Name: "Slow", Status: model.StatusYes}
	}, RunOptions{Client: &http.Client{}, IPVersion: "ipv4"})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline error, got %v", err)
	}
	if result.Name != "Slow" || result.Status != model.StatusTimeout {
		t.Fatalf("unexpected deadline result: %#v", result)
	}
}
