package in

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// JioCinema
// apis-jiocinema.voot.com 双栈 get 请求
func JioCinema(request *gorequest.SuperAgent) model.Result {
	name := "Jio Cinema"
	if request == nil {
		return model.Result{Name: name}
	}
	resp1, body1, errs1 := request.Get("https://apis-jiocinema.voot.com/location").
		Set("Accept", "application/json, text/plain, */*").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		Set("Cache-Control", "no-cache").
		Set("Connection", "keep-alive").
		Set("Origin", "https://www.jiocinema.com").
		Set("Pragma", "no-cache").
		Set("Referer", "https://www.jiocinema.com/").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "cross-site").
		Set("sec-ch-ua", "Your-UA-SecCHUA").
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()

	isBlocked1 := strings.Contains(body1, "Access Denied")
	isOK1 := strings.Contains(body1, "Success")

	resp2, body2, errs2 := request.Options("https://content-jiovoot.voot.com/psapi/voot/v1/voot-web//view/show/3500210?subNavId=38fa57ba_1706064514668&excludeTray=player-tray,subnav&responseType=common&devicePlatformType=desktop&page=1&layoutCohort=default&supportedChips=comingsoon").
		Set("Accept", "*/*").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		Set("Access-Control-Request-Headers", "app-version").
		Set("Access-Control-Request-Method", "GET").
		Set("Connection", "keep-alive").
		Set("Origin", "https://www.jiocinema.com").
		Set("Referer", "https://www.jiocinema.com/").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "cross-site").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0").
		End()

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
