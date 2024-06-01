package de

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ZDF
// ssl.zdf.de 仅 ipv4 且 get 请求
func ZDF(request *gorequest.SuperAgent) model.Result {
	name := "ZDF"
	url := "https://ssl.zdf.de/geo/de/geo.txt"
	request = request.Set("User-Agent", model.UA_Dalvik)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 000 || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ssl.zdf.de failed with code: %d", resp.StatusCode)}
}