package us

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// HBOMax
// www.hbomax.com 仅 ipv4 且 get 请求 (重定向至于 www.max.com 了)
// www.hbonow.com 仅 ipv4 且 get 请求 (重定向至于 www.max.com 了)
func HBOMax(request *gorequest.SuperAgent) model.Result {
	name := "HBO Max"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.max.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "geo-availability") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	t := strings.Split(resp.Header.Get("location"), "/")
	region := ""
	if len(t) >= 4 {
		region = strings.Split(resp.Header.Get("location"), "/")[3]
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToUpper(region)}
}
