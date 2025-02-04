package kr

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Watcha
// watcha.com 仅 ipv4 且 get 请求
func Watcha(c *http.Client) model.Result {
	name := "WATCHA"
	hostname := "watcha.com"
	if c == nil {
		return model.Result{Name: name}
	}
	// 首先检查 API 接口
	apiURL := "https://watcha.com/api/aio_browses/tvod/all?size=3"
	client := utils.Req(c)
	// 检查 API 接口
	resp1, err := client.R().Get(apiURL)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp1.Body.Close()
	if resp1.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// 检查主页面
	url := "https://watcha.com/browse/theater"
	headers := map[string]string{
		"User-Agent":                model.UA_Browser,
		"host":                      "watcha.com",
		"connection":                "keep-alive",
		"sec-ch-ua":                 model.UA_SecCHUA,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"upgrade-insecure-requests": "1",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	}
	client = utils.SetReqHeaders(client, headers)
	resp2, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	// 检查各种状态码
	if resp2.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp2.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	} else if resp2.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: "kr"}
	} else if resp2.StatusCode == 302 {
		location := resp2.Header.Get("Location")
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		switch location {
		case "/ja-JP/browse/theater":
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: "jp"}
		case "/ko-KR/browse/theater":
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: "kr"}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get watcha.com failed with code: %d", resp2.StatusCode)}
}
