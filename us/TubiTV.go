package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// TubiTV
// tubitv.com 双栈 get 请求
func TubiTV(request *gorequest.SuperAgent) model.Result {
	name := "Tubi TV"
	url := "https://tubitv.com/home"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if resp.StatusCode == 302 {
		resp2, body2, errs2 := request.Get("https://gdpr.tubi.tv").Retry(2, 5).End()
		if len(errs2) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
		}
		defer resp2.Body.Close()
		if strings.Contains(body2, "Unfortunately") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get tubitv.com failed with code: %d", resp.StatusCode)}
}