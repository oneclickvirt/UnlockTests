package vn

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// TV360
// api-v2.tv360.vn 仅 ipv4 且 get 请求 有问题
// {"errorCode":412,"message":"Xác thực hết hạn. Vui lòng đăng nhập lại1.(H-128)","data":null}
// 登录认证已过期
func TV360(c *http.Client) model.Result {
	name := "TV360"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "http://api-v2.tv360.vn/public/v1/composite/get-link?childId=998335&device_type=WEB_IPHONE&id=19474&network_device_id=prIUMaumjI7dNWKSUxFkEViFygs%3D&t=1686572228&type=film"
	userAgent := "TV360/31 CFNetwork/1402.0.8 Darwin/22.2.0"
	authorization := "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxODI1NTEzNDMiLCJ1c2VySWQiOjE4MjU1MTM0MywicHJvZmlsZUlkIjoxODI3MzM0NTUsImR2aSI6MjY5NDY3MTUzLCJjb250ZW50RmlsdGVyIjoiMTAwIiwiZ25hbWUiOiIiLCJpYXQiOjE2ODY1NzIyMDEsImV4cCI6MTY4NzE3NzAwMX0.oi0BKvATgBzPEkqR_liBrvMKXBUiWzp2BQme-biDnwiVhuta0qn_aZo6z3azLdjW5kH6PfEwEkc4K9jCfAK5rw"
	headers := map[string]string{
		"User-Agent":    userAgent,
		"userid":        "182551343",
		"devicetype":    "WEB_IPHONE",
		"deviceName":    "iPad Air 5th Gen (WiFi)",
		"profileid":     "182733455",
		"s":             "cSkV/vwUfX6tahDwe6xh9Bl0yhNs/TdWTaOJiWDt3gHekijGnNYh9i4YaUmdfBfI4oKOwvioKJ7PuKMH7ctWA6rEHeGXH/nUYOY1g7l4Umh6zoed5bBwWCgUuh5eMqdNNoptwaeCee58USTteOkbHQ==",
		"deviceid":      "69FFABD6-F9D8-4C2E-8C44-7195CF0A2930",
		"devicedrmid":   "prIUMaumjI7dNWKSUxFkEViFygs=",
		"Authorization": authorization,
		"osappversion":  "1.9.27",
		"sessionid":     "C5017358-5327-4185-999A-CA3291CC66AC",
		"zoneid":        "1",
		"Accept":        "application/json, text/html",
		"Content-Type":  "application/json",
		"osapptype":     "IPAD",
		"tv360transid":  "1686572228_69FFABD6-F9D8-4C2E-8C44-7195CF0A2930",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api-v2.tv360.vn failed with code: %d", resp.StatusCode)}
}
