package jp

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Hulu
// www.hulu.jp 或 id.hulu.jp 仅 ipv4 且 get 请求
// https://www.hulu.jp/login
func Hulu(request *gorequest.SuperAgent) model.Result {
	name := "Hulu Japan"
	if request == nil {
		return model.Result{Name: name}
	}
	request = request.Set("User-Agent", model.UA_Browser).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,"+
			"image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Set("Accept-Encoding", "gzip, deflate, br").
		Set("Cache-Control", "no-cache").
		Set("DNT", "1").
		Set("Pragma", "no-cache").
		Set("Sec-CH-UA", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`).
		Set("Sec-CH-UA-Mobile", "?0").
		Set("Sec-CH-UA-Platform", "Windows").
		Set("Sec-Fetch-Dest", "document").
		Set("Sec-Fetch-Mode", "navigate").
		Set("Sec-Fetch-Site", "none").
		Set("Sec-Fetch-User", "?1").
		Set("Upgrade-Insecure-Requests", "1")
	resp, _, errs := request.Get("https://id.hulu.jp").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.Request.URL.Path == "/restrict.html" || resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
