package africa

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Showmax
// www.showmax.com 双栈 且 get 请求
func Showmax(request *gorequest.SuperAgent) model.Result {
	name := "Showmax"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.showmax.com/"
	resp, body, errs := request.Get(url).
		Set("Host", "www.showmax.com").
		Set("Connection", "keep-alive").
		Set("Sec-Ch-UA", `"Chromium";v="124", "Microsoft Edge";v="124", "Not-A.Brand";v="99"`).
		Set("Sec-Ch-UA-Mobile", "?0").
		Set("Sec-Ch-UA-Platform", `"Windows"`).
		Set("Upgrade-Insecure-Requests", "1").
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0").
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Set("Sec-Fetch-Site", "none").
		Set("Sec-Fetch-Mode", "navigate").
		Set("Sec-Fetch-User", "?1").
		Set("Sec-Fetch-Dest", "document").
		Set("Accept-Language", "zh-CN,zh;q=0.9").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	regionStart := strings.Index(body, "activeTerritory")
	if regionStart == -1 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	regionEnd := strings.Index(body[regionStart:], "\n")
	region := strings.TrimSpace(body[regionStart+len("activeTerritory")+1 : regionStart+regionEnd])
	if region != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.showmax.com failed with code: %d", resp.StatusCode)}
}
