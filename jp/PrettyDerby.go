package jp

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// PrettyDerby
// api-umamusume.cygames.jp 双栈 且 get 请求
// 有问题 stream error: stream ID 1; INTERNAL_ERROR; received from peer
func PrettyDerby(request *gorequest.SuperAgent) model.Result {
	name := "Pretty Derby Japan"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api-umamusume.cygames.jp/"
	request = request.Set("User-Agent", model.UA_Dalvik)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api-umamusume.cygames.jp failed with code: %d", resp.StatusCode)}
}
