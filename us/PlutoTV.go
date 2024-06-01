package us

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// PlutoTV
// pluto.tv 仅 ipv4 且 get 请求
func PlutoTV(request *gorequest.SuperAgent) model.Result {
	name := "Pluto TV"
	url := "https://pluto.tv/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "thanks-for-watching") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	//if resp.StatusCode == 429 {
	//	return model.Result{Name: name, Status: model.StatusUnexpected + " Rate Limit"}
	//}
	return model.Result{Name: name, Status: model.StatusYes}
}