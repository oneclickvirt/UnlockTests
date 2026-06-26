package executor

import (
	"net/http"
	"strings"
	"syscall"
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
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

func TestShowResultUnknownStatusIsNotBlank(t *testing.T) {
	got := ShowResult(&model.Result{Name: "Test", Status: "custom failure"})
	if strings.TrimSpace(got) == "" {
		t.Fatalf("expected unknown status to be visible")
	}
	if !strings.Contains(got, "Unknown") || !strings.Contains(got, "custom failure") {
		t.Fatalf("expected unknown status details, got %q", got)
	}
}

func TestFormarPrintKeepsCDNStatusVisible(t *testing.T) {
	oldNames := Names
	oldResults := R
	defer func() {
		Names = oldNames
		R = oldResults
	}()

	Names = []string{"Netflix CDN"}
	R = []*model.Result{{Name: "Netflix CDN", Status: model.StatusYes, Region: "gb"}}

	got := FormarPrint("All")
	if !strings.Contains(got, "Netflix CDN") ||
		!strings.Contains(got, "YES") ||
		!strings.Contains(got, "Region: GB") {
		t.Fatalf("expected CDN output to include explicit status and region, got %q", got)
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

func TestResultCacheKeyIncludesTransportIdentity(t *testing.T) {
	a := &http.Client{Transport: &http.Transport{}}
	b := &http.Client{Transport: &http.Transport{}}
	keyA := resultCacheKey("Test", "ipv4", a)
	keyB := resultCacheKey("Test", "ipv4", b)
	if keyA == keyB {
		t.Fatalf("expected different cache keys for different transports")
	}
	if keyA == resultCacheKey("Test", "ipv6", a) {
		t.Fatalf("expected cache key to include IP version")
	}
}

func TestFirstDNSServerDialAddressNormalizesHostPort(t *testing.T) {
	tests := map[string]string{
		"1.1.1.1":                             "1.1.1.1:53",
		"1.1.1.1:5353 8.8.8.8:53":             "1.1.1.1:5353",
		"2606:4700:4700::1111":                "[2606:4700:4700::1111]:53",
		"[2606:4700:4700::1111]:5353,8.8.8.8": "[2606:4700:4700::1111]:5353",
	}
	for input, want := range tests {
		if got := firstDNSServerDialAddress(input); got != want {
			t.Fatalf("firstDNSServerDialAddress(%q) = %q, want %q", input, got, want)
		}
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

func TestSetupConcurrencyZeroClearsSemaphore(t *testing.T) {
	SetupConcurrency(1)
	if sem == nil {
		t.Fatalf("expected semaphore to be set")
	}
	SetupConcurrency(0)
	if sem != nil {
		t.Fatalf("expected semaphore to be cleared")
	}
}

func TestSetupDnsServersEmptyClearsResolver(t *testing.T) {
	SetupDnsServers("1.1.1.1")
	if utils.Dialer.Resolver == nil {
		t.Fatalf("expected custom resolver to be set")
	}
	SetupDnsServers("")
	if utils.Dialer.Resolver != nil {
		t.Fatalf("expected custom resolver to be cleared")
	}
}

func TestSetupInterfaceEmptyClearsBinding(t *testing.T) {
	SetupInterface("127.0.0.1")
	if utils.Dialer.LocalAddr == nil {
		t.Fatalf("expected local address to be set")
	}
	utils.Dialer.Control = func(network, address string, c syscall.RawConn) error { return nil }
	SetupInterface("")
	if utils.Dialer.LocalAddr != nil || utils.Dialer.Control != nil {
		t.Fatalf("expected interface binding to be cleared")
	}
}

func TestSetupSocksProxyEmptyRestoresHTTPProxy(t *testing.T) {
	restore := snapshotNetworkState()
	defer restore()

	SetupHttpProxy("http://127.0.0.1:8080")
	request, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	httpProxyURL, err := utils.AutoTransport.Proxy(request)
	if err != nil {
		t.Fatalf("resolve http proxy: %v", err)
	}
	if httpProxyURL == nil {
		t.Fatalf("expected HTTP proxy to be set")
	}

	SetupSocksProxy("socks5://127.0.0.1:1080")
	if utils.AutoTransport.Proxy != nil || utils.Ipv4Transport.Proxy != nil || utils.Ipv6Transport.Proxy != nil {
		t.Fatalf("expected HTTP proxy callbacks to be disabled while SOCKS proxy is active")
	}

	SetupSocksProxy("")
	restoredProxyURL, err := utils.AutoTransport.Proxy(request)
	if err != nil {
		t.Fatalf("resolve restored proxy: %v", err)
	}
	if restoredProxyURL == nil || restoredProxyURL.String() != httpProxyURL.String() {
		t.Fatalf("expected HTTP proxy to be restored, got %v want %v", restoredProxyURL, httpProxyURL)
	}
}

func snapshotNetworkState() func() {
	clientProxy := utils.ClientProxy
	autoProxy := utils.AutoTransport.Proxy
	ipv4Proxy := utils.Ipv4Transport.Proxy
	ipv6Proxy := utils.Ipv6Transport.Proxy
	autoDial := utils.AutoTransport.DialContext
	ipv4Dial := utils.Ipv4Transport.DialContext
	ipv6Dial := utils.Ipv6Transport.DialContext
	localAddr := utils.Dialer.LocalAddr
	control := utils.Dialer.Control
	resolver := utils.Dialer.Resolver
	return func() {
		utils.ClientProxy = clientProxy
		utils.AutoTransport.Proxy = autoProxy
		utils.Ipv4Transport.Proxy = ipv4Proxy
		utils.Ipv6Transport.Proxy = ipv6Proxy
		utils.AutoTransport.DialContext = autoDial
		utils.Ipv4Transport.DialContext = ipv4Dial
		utils.Ipv6Transport.DialContext = ipv6Dial
		utils.Dialer.LocalAddr = localAddr
		utils.Dialer.Control = control
		utils.Dialer.Resolver = resolver
		ClearCache()
	}
}
