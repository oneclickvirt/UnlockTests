package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

const (
	xURL      = "https://x.com/"
	xTraceURL = "https://x.com/cdn-cgi/trace"
)

var xRestrictedCountries = []string{"cn", "ir", "mm", "kp", "ru", "tm"}

func X(c *http.Client) model.Result {
	return checkX(c, xURL, xTraceURL, "x.com")
}

func checkX(c *http.Client, url, traceURL, hostname string) model.Result {
	const name = "X (formerly Twitter)"
	if c == nil {
		return model.Result{Name: name}
	}
	loc, traceStatus, ok := cloudflareTraceLocationStatus(c, hostname, traceURL)
	if !ok {
		if traceStatus == http.StatusTooManyRequests {
			return model.Result{Name: name, Status: model.StatusRateLimited, Info: "trace HTTP 429"}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("can not determine region")}
	}

	resp, err := utils.Req(c).R().
		SetHeader("User-Agent", model.UA_Browser).
		Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return model.Result{Name: name, Status: model.StatusRateLimited, Region: loc, Info: "HTTP 429"}
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnavailableForLegalReasons ||
		utils.GetRegion(loc, xRestrictedCountries) {
		return model.Result{Name: name, Status: model.StatusNo, Region: loc}
	}
	for _, keyword := range defaultAIWAFKeywords() {
		if strings.Contains(strings.ToLower(string(body)), strings.ToLower(keyword)) {
			return model.Result{Name: name, Status: model.StatusBanned, Region: loc, Info: "WAF"}
		}
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return model.Result{
			Name:   name,
			Status: model.StatusUnexpected,
			Region: loc,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
	}
	if loc == "t1" {
		loc = "tor"
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: loc}
}
