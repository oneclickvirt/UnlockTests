package executor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

const DefaultStructuredConcurrency = 20

type RunOptions struct {
	Selection    string
	IPVersion    string
	Client       *http.Client
	Concurrency  int
	UseCache     bool
	IncludeHeads bool
}

type StructuredResult struct {
	Name       string
	Status     string
	Region     string
	Info       string
	UnlockType string
	Error      string
}

func RunSelection(ctx context.Context, client *http.Client, selection, ipVersion string) ([]StructuredResult, error) {
	return RunStructured(ctx, RunOptions{
		Selection: selection,
		IPVersion: ipVersion,
		Client:    client,
	})
}

func RunStructured(ctx context.Context, opts RunOptions) ([]StructuredResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts.Selection == "" {
		opts.Selection = "0"
	}
	ipVersion, err := normalizeIPVersion(opts.IPVersion)
	if err != nil {
		return nil, err
	}
	opts.IPVersion = ipVersion
	if opts.Client == nil {
		switch opts.IPVersion {
		case "ipv4":
			opts.Client = utils.Ipv4HttpClient
		case "ipv6":
			opts.Client = utils.Ipv6HttpClient
		default:
			opts.Client = utils.AutoHttpClient
		}
	}

	runTestsMutex.Lock()
	defer runTestsMutex.Unlock()

	funcs, _, err := functionsForSelectionLocked(opts.Selection)
	if err != nil {
		return nil, err
	}
	utils.SetDNSIPVersion(opts.IPVersion)
	defer utils.SetDNSIPVersion("")
	return runFunctionsStructured(ctx, funcs, opts)
}

func normalizeIPVersion(ipVersion string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(ipVersion)) {
	case "", "4", "v4", "ip4", "ipv4":
		return "ipv4", nil
	case "6", "v6", "ip6", "ipv6":
		return "ipv6", nil
	case "0", "auto", "both", "dual", "dualstack", "ip", "any":
		return "auto", nil
	default:
		return "", fmt.Errorf("invalid IPVersion: %q", ipVersion)
	}
}

func ListPlatforms(selection string) ([]string, error) {
	if selection == "" {
		selection = "0"
	}
	runTestsMutex.Lock()
	defer runTestsMutex.Unlock()
	_, names, err := functionsForSelectionLocked(selection)
	if err != nil {
		return nil, err
	}
	return names, nil
}

type selectionState struct {
	M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT bool
	names                                           []string
	results                                         []*model.Result
}

func snapshotSelectionState() selectionState {
	return selectionState{
		M: M, TW: TW, HK: HK, JP: JP, KR: KR, NA: NA, SA: SA, EU: EU, AFR: AFR, OCEA: OCEA, SPORT: SPORT,
		names:   append([]string(nil), Names...),
		results: append([]*model.Result(nil), R...),
	}
}

func restoreSelectionState(state selectionState) {
	M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT = state.M, state.TW, state.HK, state.JP, state.KR, state.NA, state.SA, state.EU, state.AFR, state.OCEA, state.SPORT
	Names = append([]string(nil), state.names...)
	R = append([]*model.Result(nil), state.results...)
}

func functionsForSelectionLocked(selection string) ([]func(c *http.Client) model.Result, []string, error) {
	state := snapshotSelectionState()
	defer restoreSelectionState(state)

	Names = nil
	R = nil
	if !parseSelection(selection) {
		return nil, nil, fmt.Errorf("invalid selection: %q", selection)
	}
	funcs := sortedFuncList(uniqueFuncList(getFuncList()))
	names := namesFromFunctions(funcs)
	return funcs, names, nil
}

func runFunctionsStructured(ctx context.Context, funcs []func(c *http.Client) model.Result, opts RunOptions) ([]StructuredResult, error) {
	if len(funcs) == 0 {
		return []StructuredResult{}, ctx.Err()
	}
	concurrency := opts.Concurrency
	if concurrency <= 0 {
		concurrency = DefaultStructuredConcurrency
	}
	if concurrency > len(funcs) {
		concurrency = len(funcs)
	}

	results := make([]StructuredResult, len(funcs))
	jobs := make(chan int)
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	setErr := func(err error) {
		if err != nil {
			errOnce.Do(func() { firstErr = err })
		}
	}

	for worker := 0; worker < concurrency; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				if err := ctx.Err(); err != nil {
					results[idx] = structuredFromResult(contextResult(funcs[idx], err))
					setErr(err)
					continue
				}
				result, err := runFunctionWithContext(ctx, funcs[idx], opts)
				results[idx] = structuredFromResult(result)
				setErr(err)
			}
		}()
	}

sendLoop:
	for idx := range funcs {
		select {
		case <-ctx.Done():
			setErr(ctx.Err())
			for fill := idx; fill < len(funcs); fill++ {
				results[fill] = structuredFromResult(contextResult(funcs[fill], ctx.Err()))
			}
			break sendLoop
		case jobs <- idx:
		}
	}
	close(jobs)
	wg.Wait()
	return filterStructuredResults(results, opts.IncludeHeads), firstErr
}

func runFunctionWithContext(ctx context.Context, f func(c *http.Client) model.Result, opts RunOptions) (model.Result, error) {
	testInfo := safeTestInfo(f)
	testName := testInfo.Name
	if err := ctx.Err(); err != nil {
		return contextResult(f, err), err
	}
	if opts.UseCache {
		cacheKey := resultCacheKey(testName, opts.IPVersion, opts.Client)
		cacheMutex.RLock()
		if cachedResult, exists := resultCache[cacheKey]; exists {
			cacheMutex.RUnlock()
			return cachedResult, nil
		}
		cacheMutex.RUnlock()
	}

	client := clientWithContextDeadline(opts.Client, ctx)
	var result model.Result
	func() {
		defer func() {
			if r := recover(); r != nil {
				result = model.Result{Name: testName, Status: model.StatusErr, Err: fmt.Errorf("panic recovered: %v", r)}
			}
		}()
		result = utils.NormalizeResult(client, f(client), testName)
	}()
	if err := ctx.Err(); err != nil {
		return contextResult(f, err), err
	}
	if opts.UseCache {
		cacheKey := resultCacheKey(testName, opts.IPVersion, opts.Client)
		cacheMutex.Lock()
		resultCache[cacheKey] = result
		cacheMutex.Unlock()
	}
	return result, nil
}

func clientWithContextDeadline(client *http.Client, ctx context.Context) *http.Client {
	if client == nil || ctx == nil {
		return client
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		return client
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		remaining = time.Nanosecond
	}
	if client.Timeout > 0 && client.Timeout <= remaining {
		return client
	}
	clone := *client
	clone.Timeout = remaining
	return &clone
}

func safeTestInfo(f func(c *http.Client) model.Result) (result model.Result) {
	defer func() {
		if r := recover(); r != nil {
			result = model.Result{Name: "Unknown", Status: model.StatusErr, Err: fmt.Errorf("panic recovered: %v", r)}
		}
	}()
	return utils.NormalizeResult(nil, f(nil), "Unknown")
}

func contextResult(f func(c *http.Client) model.Result, err error) model.Result {
	info := safeTestInfo(f)
	status := model.StatusErr
	if errors.Is(err, context.DeadlineExceeded) {
		status = model.StatusTimeout
	} else if errors.Is(err, context.Canceled) {
		status = model.StatusErr
	}
	return model.Result{Name: info.Name, Status: status, Err: err}
}

func structuredFromResult(result model.Result) StructuredResult {
	structured := StructuredResult{
		Name:       result.Name,
		Status:     result.Status,
		Region:     result.Region,
		Info:       result.Info,
		UnlockType: result.UnlockType,
	}
	if result.Err != nil {
		structured.Error = result.Err.Error()
	}
	return structured
}

func filterStructuredResults(results []StructuredResult, includeHeads bool) []StructuredResult {
	filtered := make([]StructuredResult, 0, len(results))
	for _, result := range results {
		if !includeHeads && result.Status == model.PrintHead {
			continue
		}
		filtered = append(filtered, result)
	}
	return filtered
}
