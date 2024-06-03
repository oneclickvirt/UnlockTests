package jp

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NETRIDE
// trial.net-ride.com 双栈 get 请求
func NETRIDE(request *gorequest.SuperAgent) model.Result {
	name := "NETRIDE"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "http://trial.net-ride.com/free/free_dl.php?R_sm_code=456&R_km_url=cabb"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 || strings.Contains(body, "302 Found") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get trial.net-ride.com failed with code: %d", resp.StatusCode)}
}
