package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// MetaAI
// www.meta.ai 双栈 且 get 请求
func MetaAI(c *http.Client) model.Result {
	name := "MetaAI"
	hostname := "www.meta.ai"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.meta.ai/ajax"
	headers := map[string]string{
		"User-Agent":                model.UA_Browser,
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "en-US,en;q=0.9",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusNo, UnlockType: unlockType}
	}
	if statusCode == 400 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	urlFallback := "https://www.meta.ai/"
	respFallback, err := client.R().Get(urlFallback)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("fallback request failed: %w", err)}
	}
	defer respFallback.Body.Close()
	b, err := io.ReadAll(respFallback.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if strings.Contains(body, "GeoBlockedErrorRoot") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "GeoBlocked"}
	}
	if strings.Contains(body, "AbraHomeRoot.react") || strings.Contains(body, "AbraHomeRootConversationQuery") ||
		strings.Contains(body, "HomeRootQuery") || strings.Contains(body, "AbraRateLimitedErrorRoot") ||
		strings.Contains(body, "KadabraRootContainer") {
		var region, code string
		code = utils.ReParse(body, `"code"\s*:\s*"(.*?)"`)
		if code != "" && strings.Contains(code, "_") {
			parts := strings.Split(code, "_")
			if len(parts) >= 2 {
				region = parts[1]
			}
		} else if code != "" && len(code) < 10 {
			region = code
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		if region != "" {
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: region}
		} else {
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		}
	}
	if respFallback.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("unexpected response: ajax status=%d, home status=%d", statusCode, respFallback.StatusCode)}
}
