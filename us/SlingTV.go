package us

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// SlingTV
// p-geo.movetv.com returns Sling's own geo decision for the current IP.
func SlingTV(c *http.Client) model.Result {
	name := "Sling TV"
	hostname := "p-geo.movetv.com"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"User-Agent":      model.UA_Browser,
		"Accept":          "application/json,text/plain,*/*",
		"Origin":          "https://www.sling.com",
		"Referer":         "https://www.sling.com/",
		"Accept-Language": "en-US,en;q=0.9",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get("https://p-geo.movetv.com/geo")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnavailableForLegalReasons {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode != http.StatusOK {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get p-geo.movetv.com failed with code: %d", resp.StatusCode)}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	region, blocked := parseSlingGeoResponse(b)
	if blocked {
		return model.Result{Name: name, Status: model.StatusNo, Region: region, Info: "VPN Blocked"}
	}
	if region == "us" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
	}
	if region != "" {
		return model.Result{Name: name, Status: model.StatusNo, Region: region}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("can not parse geo response")}
}

func parseSlingGeoResponse(body []byte) (string, bool) {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		text := strings.ToLower(string(body))
		return "", strings.Contains(text, "blocked") || strings.Contains(text, "restricted")
	}
	region := firstStringValue(raw, "country", "countryCode", "country_code", "countryCodeAlpha2", "country_code_alpha2")
	blocked := firstBoolValue(raw, "ip_restricted", "isRestricted", "restricted", "blocked", "isBlocked", "vpn", "isVpn", "blacklisted")
	return normalizeSlingCountryCode(region), blocked
}

func normalizeSlingCountryCode(region string) string {
	region = strings.ToLower(strings.TrimSpace(region))
	if len(region) == 3 {
		if alpha2 := utils.ThreeToTwoCode(region); alpha2 != "" {
			return alpha2
		}
	}
	return region
}

func firstStringValue(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := raw[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstBoolValue(raw map[string]any, keys ...string) bool {
	for _, key := range keys {
		if value, ok := raw[key].(bool); ok {
			return value
		}
		if value, ok := raw[key].(string); ok {
			switch strings.ToLower(strings.TrimSpace(value)) {
			case "true", "1", "yes", "blocked", "restricted":
				return true
			}
		}
	}
	return false
}
