package vn

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// TV360
// api-v2.tv360.vn 仅 ipv4 且 get 请求 有问题
// {"errorCode":412,"message":"Xác thực hết hạn. Vui lòng đăng nhập lại1.(H-128)","data":null}
// 登录认证已过期
func TV360(request *gorequest.SuperAgent) model.Result {
	name := "TV360"
	url := "http://api-v2.tv360.vn/public/v1/composite/get-link?childId=998335&device_type=WEB_IPHONE&id=19474&network_device_id=prIUMaumjI7dNWKSUxFkEViFygs%3D&t=1686572228&type=film"
	userAgent := "TV360/31 CFNetwork/1402.0.8 Darwin/22.2.0"
	authorization := "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiIxODI1NTEzNDMiLCJ1c2VySWQiOjE4MjU1MTM0MywicHJvZmlsZUlkIjoxODI3MzM0NTUsImR2aSI6MjY5NDY3MTUzLCJjb250ZW50RmlsdGVyIjoiMTAwIiwiZ25hbWUiOiIiLCJpYXQiOjE2ODY1NzIyMDEsImV4cCI6MTY4NzE3NzAwMX0.oi0BKvATgBzPEkqR_liBrvMKXBUiWzp2BQme-biDnwiVhuta0qn_aZo6z3azLdjW5kH6PfEwEkc4K9jCfAK5rw"
	resp, body, errs := request.Get(url).
		Set("User-Agent", userAgent).
		Set("userid", "182551343").
		Set("devicetype", "WEB_IPHONE").
		Set("deviceName", "iPad Air 5th Gen (WiFi)").
		Set("profileid", "182733455").
		Set("s", "cSkV/vwUfX6tahDwe6xh9Bl0yhNs/TdWTaOJiWDt3gHekijGnNYh9i4YaUmdfBfI4oKOwvioKJ7PuKMH7ctWA6rEHeGXH/nUYOY1g7l4Umh6zoed5bBwWCgUuh5eMqdNNoptwaeCee58USTteOkbHQ==").
		Set("deviceid", "69FFABD6-F9D8-4C2E-8C44-7195CF0A2930").
		Set("devicedrmid", "prIUMaumjI7dNWKSUxFkEViFygs=").
		Set("Authorization", authorization).
		Set("osappversion", "1.9.27").
		Set("sessionid", "C5017358-5327-4185-999A-CA3291CC66AC").
		Set("zoneid", "1").
		Set("Accept", "application/json, text/html").
		Set("Content-Type", "application/json").
		Set("osapptype", "IPAD").
		Set("tv360transid", "1686572228_69FFABD6-F9D8-4C2E-8C44-7195CF0A2930").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api-v2.tv360.vn failed with code: %d", resp.StatusCode)}
}
