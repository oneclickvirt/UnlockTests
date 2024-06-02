package us

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// PeacockTV
// www.peacocktv.com 双栈 get 请求
func PeacockTV(request *gorequest.SuperAgent) model.Result {
	name := "Peacock TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.peacocktv.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(resp.Header.Get("location"), "unavailable") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
