package kr

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"io"
	"strings"
)

// SPOTVNOW
// edge.api.brightcove.com 仅 ipv4 且 get 请求
func SPOTVNOW(request *gorequest.SuperAgent) model.Result {
	name := "SPOTV NOW"
	url := "https://edge.api.brightcove.com/playback/v1/accounts/5764318566001/videos/6349973203112"
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("User-Agnet", model.UA_Browser)
	client.Headers.Set("sec-ch-ua", model.UA_SecCHUA)
	client.Headers.Set("referer", "https://www.spotvnow.co.kr/")
	client.Headers.Set("origin", "https://www.spotvnow.co.kr")
	client.Headers.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8")
	client.Headers.Set("sec-ch-ua-mobile", "?0")
	client.Headers.Set("sec-ch-ua-platform", "\"Windows\"")
	client.Headers.Set("sec-fetch-dest", "empty")
	client.Headers.Set("sec-fetch-mode", "cors")
	client.Headers.Set("sec-fetch-site", "cross-site")
	client.Headers.Set("accept", "application/json;pk=BCpkADawqM0U3mi_PT566m5lvtapzMq3Uy7ICGGjGB6v4Ske7ZX_ynzj8ePedQJhH36nym_5mbvSYeyyHOOdUsZovyg2XlhV6rRspyYPw_USVNLaR0fB_AAL2HSQlfuetIPiEzbUs1tpNF9NtQxt3BAPvXdOAsvy1ltLPWMVzJHiw9slpLRgI2NUufc")
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	var res1 struct {
		ErrorSubcode string `json:"error_subcode"`
		AccountId    string `json:"account_id"`
	}
	var res2 []struct {
		ClientGeo    string `json:"client_geo"`
		ErrorSubcode string `json:"error_subcode"`
		ErrorCode    string `json:"error_code"`
		Message      string `json:"message"`
	}
	if err := json.Unmarshal(b, &res1); err != nil {
		if err := json.Unmarshal(b, &res2); err != nil {
			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res2[0].ErrorSubcode == "CLIENT_GEO" {
			return model.Result{Name: name, Status: model.StatusNo, Region: res2[0].ClientGeo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res1.AccountId != "0" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "kr"}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get edge.api.brightcove.com with code: %d", resp.StatusCode)}
}
