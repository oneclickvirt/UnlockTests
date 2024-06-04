package in

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// JioCinema
// apis-jiocinema.voot.com 双栈 get 请求
func JioCinema(c *http.Client) model.Result {
	name := "Jio Cinema"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Cache-Control":      "no-cache",
		"Connection":         "keep-alive",
		"Origin":             "https://www.jiocinema.com",
		"Pragma":             "no-cache",
		"Referer":            "https://www.jiocinema.com/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	url1 := "https://apis-jiocinema.voot.com/location"
	resp1, body1, errs1 := request.Get(url1).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()

	isBlocked1 := strings.Contains(body1, "Access Denied")
	isOK1 := strings.Contains(body1, "Success")

	headers2 := map[string]string{
		"Accept":                         "*/*",
		"Accept-Language":                "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Access-Control-Request-Headers": "app-version",
		"Access-Control-Request-Method":  "GET",
		"Connection":                     "keep-alive",
		"Origin":                         "https://www.jiocinema.com",
		"Referer":                        "https://www.jiocinema.com/",
		"Sec-Fetch-Dest":                 "empty",
		"Sec-Fetch-Mode":                 "cors",
		"Sec-Fetch-Site":                 "cross-site",
		"User-Agent":                     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
	}
	request2 := utils.Gorequest(c)
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	url2 := "https://content-jiovoot.voot.com/psapi/voot/v1/voot-web//view/show/3500210?subNavId=38fa57ba_1706064514668&excludeTray=player-tray,subnav&responseType=common&devicePlatformType=desktop&page=1&layoutCohort=default&supportedChips=comingsoon"
	resp2, body2, errs2 := request2.Options(url2).End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()

	isBlocked2 := strings.Contains(body2, "JioCinema is unavailable at your location")
	isOK2 := strings.Contains(body2, "Ok")

	if isBlocked1 || isBlocked2 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if isOK1 && isOK2 {
		return model.Result{Name: name, Status: model.StatusYes}
	}

	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get apis-jiocinema.voot.com failed with code: %d", resp2.StatusCode)}
}
