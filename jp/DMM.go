package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// DMM
// bitcoin.dmm.com 仅 ipv4 且 get 请求
func DMM(request *gorequest.SuperAgent) model.Result {
	name := "DMM"
	url := "https://bitcoin.dmm.com"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "This page is not available in your area") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "暗号資産") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get bitcoin.dmm.com failed with code: %d", resp.StatusCode)}
}
