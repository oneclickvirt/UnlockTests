package uk

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Channel4
// www.channel4.com 仅 ipv4 且 get 请求
func Channel4(request *gorequest.SuperAgent) model.Result {
	name := "Channel 4"
	url := "https://www.channel4.com/simulcast/channels/C4"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.channel4.com failed with code: %d", resp.StatusCode)}
}
