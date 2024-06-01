package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ATTNOW - DirectvStream
// www.atttvnow.com 双栈 且 get 请求
func DirectvStream(request *gorequest.SuperAgent) model.Result {
	name := "Directv Stream"
	url := "https://www.atttvnow.com/"
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes,
		Err: fmt.Errorf("get www.atttvnow.com failed with code: %d", resp.StatusCode)}
}
