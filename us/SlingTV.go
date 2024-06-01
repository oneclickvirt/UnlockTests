package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// SlingTV
// www.sling.com 双栈 且 get 请求
func SlingTV(request *gorequest.SuperAgent) model.Result {
	name := "Sling TV"
	url := "https://www.sling.com/"
	request = request.Set("User-Agent", model.UA_Dalvik)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.sling.com failed with code: %d", resp.StatusCode)}
}
