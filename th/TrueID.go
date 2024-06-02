package th

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func getStringBetween(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	posLast := strings.Index(value[posFirstAdjusted:], b)
	if posLast == -1 {
		return ""
	}
	return value[posFirstAdjusted : posFirstAdjusted+posLast]
}

// TrueID
// tv.trueid.net 双栈 get 请求
func TrueID(request *gorequest.SuperAgent) model.Result {
	name := "TrueID"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://tv.trueid.net/th-en/live/thairathtv-hd"
	headers := map[string]string{
		"User-Agent":                "{UA_Browser}",
		"sec-ch-ua":                 "{UA_SecCHUA}",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}
	for key, value := range headers {
		request = request.Set(key, value)
	}
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	channelId := getStringBetween(body, `"channelId":"`, `"`)
	authUser := getStringBetween(body, `"buildId":"`, `"`)
	if len(authUser) < 11 {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("authUser len < 11")}
	}
	authKey := authUser[10:]
	apiURL := fmt.Sprintf("https://tv.trueid.net/api/stream/checkedPlay?channelId=%s&lang=en&country=th", channelId)
	authHeader := fmt.Sprintf("%s:%s", authUser, authKey)
	request = request.Set("Authorization", authHeader).
		Set("accept", "application/json, text/plain, */*").
		Set("referer", url)
	resp, body, errs = request.Get(apiURL).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	result := getStringBetween(body, `"billboardType":"`, `"`)
	if result == "GEO_BLOCK" || strings.Contains(body, "Access denied") || resp.StatusCode == 401 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if result == "LOADING" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get  failed with code: %d", resp.StatusCode)}
}
