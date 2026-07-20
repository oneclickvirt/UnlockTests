package executor

import (
	"net/http"
	"sort"
	"strings"
	"syscall"
	"testing"

	"github.com/mattn/go-runewidth"
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

func TestShowResultRateLimited(t *testing.T) {
	got := ShowResult(&model.Result{Name: "Test", Status: model.StatusRateLimited, Info: "HTTP 429", Region: "us"})
	for _, want := range []string{"Rate Limited", "HTTP 429", "Region: US"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected rate-limit output to contain %q, got %q", want, got)
		}
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

func TestFormarPrintAlignsWideProviderNames(t *testing.T) {
	state := snapshotSelectionState()
	defer restoreSelectionState(state)
	Names = []string{"中文平台", "ASCII"}
	R = []*model.Result{
		{Name: "中文平台", Status: model.StatusYes},
		{Name: "ASCII", Status: model.StatusYes},
	}
	got := FormarPrint("Test")
	lines := strings.Split(got, "\n")
	var chinese, ascii string
	for _, line := range lines {
		if strings.HasPrefix(line, "中文平台") {
			chinese = line
		}
		if strings.HasPrefix(line, "ASCII") {
			ascii = line
		}
	}
	if chinese == "" || ascii == "" {
		t.Fatalf("formatted output missing fixture lines: %q", got)
	}
	chineseResult := strings.Index(chinese, "YES")
	asciiResult := strings.Index(ascii, "YES")
	if chineseResult < 0 || asciiResult < 0 {
		t.Fatalf("formatted output missing result text: %q", got)
	}
	chinesePrefix := strings.SplitN(chinese, "\x1b", 2)[0]
	asciiPrefix := strings.SplitN(ascii, "\x1b", 2)[0]
	if runewidth.StringWidth(chinesePrefix) != runewidth.StringWidth(asciiPrefix) {
		t.Fatalf("provider names are not display-aligned: chinese=%q ascii=%q", chinese, ascii)
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
	if M || TW || HK || JP || KR || NA || SA || EU || AFR || OCEA || SPORT || AI {
		t.Fatalf("expected invalid selection to reset all selection flags")
	}
}

func TestParseSelectionSupportsAIAndCommaSeparatedItems(t *testing.T) {
	if !parseSelection("0,21") {
		t.Fatalf("expected comma-separated selection to parse")
	}
	if !M || !AI {
		t.Fatalf("expected global and AI selections to be enabled")
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

func TestReferenceProvidersArePresentInExpectedSections(t *testing.T) {
	oldNames := Names
	defer func() { Names = oldNames }()

	tests := map[string]struct {
		funcs []func(c *http.Client) model.Result
		names []string
	}{
		"global": {
			funcs: Multination(),
			names: []string{"Bilibili Anime", "Coze", "Dola AI", "Microsoft Copilot", "Poe", "WeTV", "X (formerly Twitter)"},
		},
		"ai": {
			funcs: AIPlatforms(),
			names: []string{"Coze", "DeepSeek", "Dola AI", "Grok", "Kimi", "Mistral AI", "Perplexity AI", "Poe"},
		},
		"europe": {
			funcs: Europe(),
			names: []string{"TNTSports"},
		},
		"hong kong": {
			funcs: HongKong(),
			names: []string{"Hoy TV"},
		},
		"india": {
			funcs: India(),
			names: []string{"Tata Play"},
		},
		"southeast asia": {
			funcs: SouthEastAsia(),
			names: []string{"Galaxy Play", "K+", "TV360", "Sooka"},
		},
	}
	for section, tt := range tests {
		got := namesFromFuncList(tt.funcs)
		for _, name := range tt.names {
			if !got[name] {
				t.Fatalf("expected %s section to include %q", section, name)
			}
		}
	}
}

func TestPlatformSectionsAreAlphabetized(t *testing.T) {
	oldNames := Names
	defer func() { Names = oldNames }()

	for section, build := range platformSectionBuilders() {
		Names = nil
		assertAlphabetized(t, section, orderedNamesFromFuncList(build()))
	}
}

func TestListPlatformsReturnsAlphabetizedUniqueNames(t *testing.T) {
	for _, selection := range []string{"0", "10", "14", "19", "20", "21", "0 10", "0 19", "0,21"} {
		names, err := ListPlatforms(selection)
		if err != nil {
			t.Fatalf("ListPlatforms(%q) returned error: %v", selection, err)
		}
		assertNoDuplicateNames(t, "selection "+selection, names)
		assertAlphabetized(t, "selection "+selection, names)
	}
}

func TestRegionalProvidersAreNotClassifiedAsGlobal(t *testing.T) {
	oldNames := Names
	defer func() { Names = oldNames }()

	Names = nil
	globalNames := namesFromFuncList(Multination())
	for _, name := range []string{"Acorn TV", "AMC+", "BritBox", "HBO Max", "HotStar", "Viaplay"} {
		if globalNames[name] {
			t.Fatalf("regional provider %q should not be in global section", name)
		}
	}
}

func TestBilibiliBrandStaysInGlobalSectionOnly(t *testing.T) {
	oldNames := Names
	defer func() { Names = oldNames }()

	Names = nil
	globalNames := namesFromFuncList(Multination())
	if !globalNames["Bilibili Anime"] {
		t.Fatalf("expected Bilibili Anime to be in global section")
	}

	for section, build := range geographicSectionBuilders() {
		Names = nil
		for _, name := range orderedNamesFromFuncList(build()) {
			if strings.Contains(strings.ToLower(name), "bilibili") {
				t.Fatalf("Bilibili provider %q should stay in global section, found in %s", name, section)
			}
		}
	}
}

func TestReferenceProviderSectionsHaveNoDuplicateNames(t *testing.T) {
	oldNames := Names
	defer func() { Names = oldNames }()

	sections := map[string][]func(c *http.Client) model.Result{
		"global":         Multination(),
		"europe":         Europe(),
		"hong kong":      HongKong(),
		"india":          India(),
		"southeast asia": SouthEastAsia(),
		"ipv6 global":    IPV6Multination(),
		"ai":             AIPlatforms(),
	}
	for section, funcs := range sections {
		seen := map[string]bool{}
		for _, f := range funcs {
			result := f(nil)
			if result.Status == model.PrintHead || result.Name == "" {
				continue
			}
			if seen[result.Name] {
				t.Fatalf("duplicate provider %q in %s section", result.Name, section)
			}
			seen[result.Name] = true
		}
	}
}

func TestFunctionsForTestNamesAcceptsCommaSeparatedNames(t *testing.T) {
	funcs, names, missing := functionsForTestNamesLocked("Coze,Poe,Perplexity AI")
	if len(missing) != 0 {
		t.Fatalf("unexpected missing tests: %v", missing)
	}
	if len(funcs) != 3 {
		t.Fatalf("expected 3 funcs, got %d", len(funcs))
	}
	got := map[string]bool{}
	for _, name := range names {
		got[name] = true
	}
	for _, name := range []string{"Coze", "Poe", "Perplexity AI"} {
		if !got[name] {
			t.Fatalf("expected selected tests to include %q, got %v", name, names)
		}
	}
}

func TestFunctionsForTestNamesMatchesFunctionNames(t *testing.T) {
	_, names, missing := functionsForTestNamesLocked("MistralAI,Kimi")
	if len(missing) != 0 {
		t.Fatalf("unexpected missing tests: %v", missing)
	}
	got := strings.Join(names, "\n")
	if !strings.Contains(got, "Mistral AI") || !strings.Contains(got, "Kimi") {
		t.Fatalf("expected function-name lookup to match providers, got %v", names)
	}
}

func TestFormatVersionedHeader(t *testing.T) {
	tests := map[string]struct {
		netType string
		title   string
		want    string
	}{
		"ipv4 chinese section": {netType: "ipv4", title: "跨国平台", want: "IPV4 跨国平台"},
		"ipv6 selected tests":  {netType: "ipv6", title: "Selected Tests", want: "IPV6 Selected Tests"},
		"empty title":          {netType: "tcp4", title: "", want: "IPV4"},
	}
	for name, tt := range tests {
		if got := formatVersionedHeader(tt.netType, tt.title); got != tt.want {
			t.Fatalf("%s: got %q, want %q", name, got, tt.want)
		}
	}
}

func platformSectionBuilders() map[string]func() []func(c *http.Client) model.Result {
	builders := userVisibleSectionBuilders()
	builders["ipv6 global"] = IPV6Multination
	return builders
}

func userVisibleSectionBuilders() map[string]func() []func(c *http.Client) model.Result {
	builders := geographicSectionBuilders()
	builders["global"] = Multination
	builders["sports"] = Sport
	builders["ai"] = AIPlatforms
	return builders
}

func geographicSectionBuilders() map[string]func() []func(c *http.Client) model.Result {
	return map[string]func() []func(c *http.Client) model.Result{
		"africa":         Africa,
		"europe":         Europe,
		"hong kong":      HongKong,
		"india":          India,
		"japan":          Japan,
		"korea":          Korea,
		"north america":  NorthAmerica,
		"oceania":        Oceania,
		"south america":  SouthAmerica,
		"southeast asia": SouthEastAsia,
		"taiwan":         Taiwan,
	}
}

func namesFromFuncList(funcs []func(c *http.Client) model.Result) map[string]bool {
	names := map[string]bool{}
	for _, f := range funcs {
		result := f(nil)
		if result.Name != "" && result.Status != model.PrintHead {
			names[result.Name] = true
		}
	}
	return names
}

func orderedNamesFromFuncList(funcs []func(c *http.Client) model.Result) []string {
	names := []string{}
	for _, f := range funcs {
		result := f(nil)
		if result.Name != "" && result.Status != model.PrintHead {
			names = append(names, result.Name)
		}
	}
	return names
}

func assertAlphabetized(t *testing.T, label string, names []string) {
	t.Helper()
	want := append([]string(nil), names...)
	sort.SliceStable(want, func(i, j int) bool {
		left := strings.ToLower(want[i])
		right := strings.ToLower(want[j])
		if left == right {
			return want[i] < want[j]
		}
		return left < right
	})
	if strings.Join(names, "\n") != strings.Join(want, "\n") {
		t.Fatalf("%s names are not alphabetized:\ngot  %v\nwant %v", label, names, want)
	}
}

func assertNoDuplicateNames(t *testing.T, label string, names []string) {
	t.Helper()
	seen := map[string]bool{}
	for _, name := range names {
		if seen[name] {
			t.Fatalf("%s has duplicate provider %q", label, name)
		}
		seen[name] = true
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
	if !strings.Contains(got, "[ IPV6 Taiwan ]") {
		t.Fatalf("expected IPv6 output to use selected platform title, got %q", got)
	}
}

func TestFinallyPrintResultIPv4PrefixesSectionTitle(t *testing.T) {
	resetOptions()
	defer func() {
		resetOptions()
		Names = nil
		R = nil
	}()
	M = true
	Names = []string{"Example"}
	R = []*model.Result{{Name: "Example", Status: model.StatusYes}}
	got := finallyPrintResult("zh", "ipv4")
	if !strings.Contains(got, "[ IPV4 跨国平台 ]") {
		t.Fatalf("expected IPv4 output to prefix platform title, got %q", got)
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

func TestDefaultTransportDoesNotReadEnvironmentProxy(t *testing.T) {
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:9")
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:9")
	SetupHttpProxy("")
	if utils.ClientProxy != nil || utils.AutoTransport.Proxy != nil || utils.Ipv4Transport.Proxy != nil || utils.Ipv6Transport.Proxy != nil {
		t.Fatal("ambient proxy environment must not affect the default transports")
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
