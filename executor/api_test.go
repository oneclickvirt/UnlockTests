package executor

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
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
	if got.Transport == source.Transport {
		t.Fatalf("expected transport to be wrapped with caller context")
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
	if got == source || got.Timeout != source.Timeout {
		t.Fatalf("expected context-bound clone with original shorter timeout, got %#v", got)
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

func TestStructuredResultPreservesRateLimitedStatus(t *testing.T) {
	got := structuredFromResult(model.Result{
		Name:   "Limited",
		Status: model.StatusRateLimited,
		Region: "us",
		Info:   "HTTP 429",
	})
	if got.Name != "Limited" || got.Status != model.StatusRateLimited || got.Region != "us" || got.Info != "HTTP 429" {
		t.Fatalf("unexpected structured result: %#v", got)
	}
}

func TestRunStructuredAutoRunsBothIPVersions(t *testing.T) {
	oldMultination := M
	defer func() { M = oldMultination }()

	// Selection validation and the dual-stack loop are exercised here with a
	// canceled context so no provider performs a live HTTP request.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results, err := RunStructured(ctx, RunOptions{Selection: "21", IPVersion: "auto", Concurrency: 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled dual-stack run, got %v", err)
	}
	seen := map[string]bool{}
	for _, result := range results {
		seen[result.IPVersion] = true
	}
	if !seen["ipv4"] || !seen["ipv6"] {
		t.Fatalf("auto run did not return both IP versions: %#v", seen)
	}
}

func TestSplitVersionContextReservesBudgetForBothStacks(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	first, firstCancel := splitVersionContext(parent, 2)
	defer firstCancel()
	firstDeadline, ok := first.Deadline()
	if !ok {
		t.Fatal("first IP version has no deadline")
	}
	remaining := time.Until(firstDeadline)
	if remaining < 800*time.Millisecond || remaining > 1200*time.Millisecond {
		t.Fatalf("first IP version budget = %s, want about half", remaining)
	}

	second, secondCancel := splitVersionContext(parent, 1)
	defer secondCancel()
	secondDeadline, ok := second.Deadline()
	if !ok || !secondDeadline.After(firstDeadline) {
		t.Fatalf("second IP version did not receive the remaining budget: first=%s second=%s", firstDeadline, secondDeadline)
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

// Non-HTTP provider code cannot be preempted without changing every provider
// signature; the executor still reports its expired context after it returns.
func TestRunFunctionWithContextDeadlineForNonHTTPProvider(t *testing.T) {
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

func TestRunFunctionWithContextCancelsInFlightProviderHTTP(t *testing.T) {
	started := make(chan struct{})
	requestCanceled := make(chan struct{})
	var startedOnce, canceledOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		startedOnce.Do(func() { close(started) })
		<-request.Context().Done()
		canceledOnce.Do(func() { close(requestCanceled) })
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	resultChannel := make(chan struct {
		result model.Result
		err    error
	}, 1)
	go func() {
		result, err := runFunctionWithContext(ctx, func(client *http.Client) model.Result {
			if client == nil {
				return model.Result{Name: "HTTP fixture"}
			}
			_, requestErr := utils.Req(client).R().Get(server.URL)
			return model.Result{Name: "HTTP fixture", Status: model.StatusNetworkErr, Err: requestErr}
		}, RunOptions{Client: server.Client(), IPVersion: "ipv4"})
		resultChannel <- struct {
			result model.Result
			err    error
		}{result: result, err: err}
	}()
	<-started
	cancel()
	select {
	case <-requestCanceled:
	case <-time.After(time.Second):
		t.Fatal("provider HTTP request did not inherit caller cancellation")
	}
	select {
	case got := <-resultChannel:
		if !errors.Is(got.err, context.Canceled) || got.result.Status != model.StatusErr {
			t.Fatalf("canceled provider result = %#v, %v", got.result, got.err)
		}
	case <-time.After(time.Second):
		t.Fatal("provider did not return after HTTP cancellation")
	}
}

func TestValidateNetworkOptions(t *testing.T) {
	if err := validateNetworkOptions(RunOptions{HTTPProxy: "http://proxy.fixture", SOCKSProxy: "socks5://proxy.fixture"}); err == nil {
		t.Fatal("expected mutually exclusive proxy options to fail")
	}
	if err := validateNetworkOptions(RunOptions{DNSServers: ",,;"}); err == nil {
		t.Fatal("expected invalid DNS options to fail")
	}
	if err := validateNetworkOptions(RunOptions{HTTPProxy: "http://proxy.fixture"}); err != nil {
		t.Fatalf("valid HTTP proxy rejected: %v", err)
	}
}

func TestRunNamedStructuredRejectsUnknownProvider(t *testing.T) {
	results, err := RunNamedStructured(context.Background(), RunOptions{IPVersion: "ipv4"}, "definitely-not-a-provider")
	if err == nil || len(results) != 0 {
		t.Fatalf("unknown named provider result = %#v, %v", results, err)
	}
}

func TestValidateNetworkOptionsRejectsUnsupportedInterfaceNameBinding(t *testing.T) {
	if runtime.GOOS == "linux" {
		if err := validateNetworkOptions(RunOptions{Interface: "eth0"}); err != nil {
			t.Fatalf("Linux interface binding rejected: %v", err)
		}
		return
	}
	if err := validateNetworkOptions(RunOptions{Interface: "en0"}); err == nil {
		t.Fatalf("interface name binding silently accepted on %s", runtime.GOOS)
	}
}

func TestStructuredNetworkOptionsRestoreGlobalState(t *testing.T) {
	oldIP := utils.GetDNSIPVersion()
	defer utils.SetDNSIPVersion(oldIP)
	if err := validateNetworkOptions(RunOptions{IPVersion: "ipv4"}); err != nil {
		t.Fatal(err)
	}
	restore := applyStructuredNetworkOptions(RunOptions{IPVersion: "ipv4", DNSServers: "192.0.2.53"})
	if got := utils.GetDNSIPVersion(); got != "ipv4" {
		t.Fatalf("DNS IP version = %q", got)
	}
	if got := utils.CustomDNSServers(); len(got) != 1 || got[0] != "192.0.2.53" {
		t.Fatalf("custom DNS state = %#v", got)
	}
	restore()
	if got := utils.GetDNSIPVersion(); got != oldIP {
		t.Fatalf("DNS IP version was not restored: %q", got)
	}
}
