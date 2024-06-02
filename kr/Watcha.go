package kr

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Watcha
// watcha.com 仅 ipv4 且 get 请求
func Watcha(request *gorequest.SuperAgent) model.Result {
	name := "WATCHA"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://watcha.com/"
	request = request.Set("User-Agent", model.UA_Browser).
		Set("host", "watcha.com").
		Set("connection", "keep-alive").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"Windows\"").
		Set("upgrade-insecure-requests", "1").
		Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' -H 'sec-fetch-site: none' -H 'sec-fetch-mode: navigate' -H 'sec-fetch-user: ?1' -H 'sec-fetch-dest: document' -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get watcha.com failed with code: %d", resp.StatusCode)}
}
