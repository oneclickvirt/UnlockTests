package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	req "github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

type aiStatusProbe struct {
	name       string
	hostname   string
	url        string
	noRedirect bool
	okCodes    map[int]bool
	noCodes    map[int]bool
}

type aiRegionalProbe struct {
	name                string
	hostname            string
	url                 string
	traceURL            string
	noRedirect          bool
	okCodes             map[int]bool
	noCodes             map[int]bool
	bannedCodes         map[int]bool
	forbiddenCodes      map[int]bool
	supportCountries    []string
	restrictedCountries []string
	wafKeywords         []string
}

var aiGlobalRestrictedCountries = []string{"cn", "ru", "ir", "kp", "cu", "sy"}

var mistralAIRestrictedCountries = []string{"ru", "by", "kp", "ir", "sy", "cu", "cn", "tm"}

func checkAIStatus(c *http.Client, probe aiStatusProbe) model.Result {
	if c == nil {
		return model.Result{Name: probe.name}
	}
	client := utils.Req(c)
	if probe.noRedirect {
		client.SetRedirectPolicy(req.NoRedirectPolicy())
	}
	resp, err := client.R().Get(probe.url)
	if err != nil {
		return utils.HandleNetworkError(c, probe.hostname, err, probe.name)
	}
	defer resp.Body.Close()
	switch {
	case probe.okCodes[resp.StatusCode]:
		return model.Result{Name: probe.name, Status: model.StatusYes}
	case probe.noCodes[resp.StatusCode]:
		return model.Result{Name: probe.name, Status: model.StatusNo}
	default:
		return model.Result{
			Name:   probe.name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
	}
}

func checkAIRegionalStatus(c *http.Client, probe aiRegionalProbe) model.Result {
	if c == nil {
		return model.Result{Name: probe.name}
	}
	client := utils.Req(c)
	if probe.noRedirect {
		client.SetRedirectPolicy(req.NoRedirectPolicy())
	}
	resp, err := client.R().
		SetHeader("User-Agent", model.UA_Browser).
		Get(probe.url)
	if err != nil {
		return utils.HandleNetworkError(c, probe.hostname, err, probe.name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, probe.hostname, err, probe.name)
	}
	bodyText := strings.ToLower(string(body))
	for _, keyword := range probe.wafKeywords {
		if strings.Contains(bodyText, strings.ToLower(keyword)) {
			return model.Result{Name: probe.name, Status: model.StatusBanned, Info: "WAF"}
		}
	}
	switch {
	case probe.forbiddenCodes[resp.StatusCode]:
		return aiForbiddenRegionResult(c, probe)
	case probe.bannedCodes[resp.StatusCode]:
		loc, ok := cloudflareTraceLocation(c, probe.hostname, probe.traceURL)
		if ok {
			return model.Result{Name: probe.name, Status: model.StatusBanned, Region: loc}
		}
		return model.Result{Name: probe.name, Status: model.StatusBanned}
	case probe.noCodes[resp.StatusCode]:
		loc, ok := cloudflareTraceLocation(c, probe.hostname, probe.traceURL)
		if ok {
			return model.Result{Name: probe.name, Status: model.StatusNo, Region: loc}
		}
		return model.Result{Name: probe.name, Status: model.StatusNo}
	case probe.okCodes[resp.StatusCode]:
		loc, ok := cloudflareTraceLocation(c, probe.hostname, probe.traceURL)
		if !ok {
			return model.Result{Name: probe.name, Status: model.StatusYes}
		}
		return aiRegionResult(probe.name, probe.hostname, loc, probe.supportCountries, probe.restrictedCountries)
	default:
		loc, ok := cloudflareTraceLocation(c, probe.hostname, probe.traceURL)
		if ok && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return aiRegionResult(probe.name, probe.hostname, loc, probe.supportCountries, probe.restrictedCountries)
		}
		return model.Result{
			Name:   probe.name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
	}
}

func cloudflareTraceLocation(c *http.Client, hostname, traceURL string) (string, bool) {
	if strings.TrimSpace(traceURL) == "" {
		traceURL = "https://" + hostname + "/cdn-cgi/trace"
	}
	resp, err := utils.Req(c).R().
		SetHeader("User-Agent", model.UA_Browser).
		Get(traceURL)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", false
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	loc := parseCloudflareTraceLocation(string(b))
	if loc == "" {
		return "", false
	}
	return loc, true
}

func parseCloudflareTraceLocation(body string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "loc=") {
			return strings.ToLower(strings.TrimSpace(strings.TrimPrefix(line, "loc=")))
		}
	}
	return ""
}

func aiRegionResult(name, hostname, loc string, supportCountries, restrictedCountries []string) model.Result {
	if loc == "t1" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "tor"}
	}
	if !aiRegionAllowed(loc, supportCountries, restrictedCountries) {
		return model.Result{Name: name, Status: model.StatusNo, Region: loc}
	}
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{Name: name, Status: model.StatusYes, Region: loc, UnlockType: unlockType}
}

func aiForbiddenRegionResult(c *http.Client, probe aiRegionalProbe) model.Result {
	loc, ok := cloudflareTraceLocation(c, probe.hostname, probe.traceURL)
	if !ok {
		return model.Result{Name: probe.name, Status: model.StatusBanned}
	}
	if aiRegionAllowed(loc, probe.supportCountries, probe.restrictedCountries) {
		return model.Result{Name: probe.name, Status: model.StatusBanned, Region: loc}
	}
	return model.Result{Name: probe.name, Status: model.StatusNo, Region: loc}
}

func aiRegionAllowed(loc string, supportCountries, restrictedCountries []string) bool {
	loc = strings.ToLower(strings.TrimSpace(loc))
	if loc == "" {
		return false
	}
	if len(restrictedCountries) > 0 && utils.GetRegion(loc, restrictedCountries) {
		return false
	}
	return len(supportCountries) == 0 || utils.GetRegion(loc, supportCountries)
}

func defaultAIWAFKeywords() []string {
	return []string{
		"attention required",
		"cf-chl",
		"checking your browser",
		"just a moment",
		"access denied",
		"request blocked",
		"unusual traffic",
	}
}
