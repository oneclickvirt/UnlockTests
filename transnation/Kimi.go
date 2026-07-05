package transnation

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

var kimiRegionRegex = regexp.MustCompile(`"useRegion":"REGION_([^"]+)"`)

func Kimi(c *http.Client) model.Result {
	name := "Kimi"
	hostname := "www.kimi.com"
	if c == nil {
		return model.Result{Name: name}
	}
	resp, err := utils.Req(c).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", model.UA_Browser).
		SetBody("{}").
		Post("https://www.kimi.com/apiv2/kimi.gateway.order.v1.GoodsService/ListGoods")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusForbidden:
		return model.Result{Name: name, Status: model.StatusNo}
	case http.StatusOK:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		body := string(b)
		for _, keyword := range defaultAIWAFKeywords() {
			if strings.Contains(strings.ToLower(body), strings.ToLower(keyword)) {
				return model.Result{Name: name, Status: model.StatusBanned, Info: "WAF"}
			}
		}
		if matches := kimiRegionRegex.FindStringSubmatch(body); len(matches) > 1 {
			region := strings.ToLower(matches[1])
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:        name,
		hostname:    hostname,
		url:         "https://www.kimi.com/",
		traceURL:    "https://www.kimi.com/cdn-cgi/trace",
		noRedirect:  true,
		okCodes:     map[int]bool{http.StatusOK: true, http.StatusMovedPermanently: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:     map[int]bool{http.StatusUnavailableForLegalReasons: true},
		bannedCodes: map[int]bool{http.StatusForbidden: true},
		wafKeywords: defaultAIWAFKeywords(),
	})
}
