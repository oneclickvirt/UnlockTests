package uk

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// ITVX
// simulcast.itv.com 仅 ipv4 且 get 请求
func ITVX(request *gorequest.SuperAgent) model.Result {
	name := "ITV Hub"
	resp, body, errs := request.Get("https://simulcast.itv.com/playlist/itvonline/ITV").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || strings.Contains(body, "Outside Of Allowed Geographic Region") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get simulcast.itv.com failed with code: %d", resp.StatusCode)}
}
