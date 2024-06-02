package uk

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// SkyGo
// skyid.sky.com 仅 ipv4 且 get 请求
func SkyGo(request *gorequest.SuperAgent) model.Result {
	name := "Sky Go"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://skyid.sky.com/authorise/skygo?response_type=token&client_id=sky&appearance=compact&redirect_uri=skygo://auth"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "You don't have permission to access") || resp.StatusCode == 403 ||
		strings.Contains(body, "Access Denied") { // || resp.StatusCode == 451
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 || resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get skyid.sky.com failed with code: %d", resp.StatusCode)}
}
