package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// PCRJP
// api-priconne-redive.cygames.jp 仅 ipv4 且 get 请求
// 有问题 stream error: stream ID 1; INTERNAL_ERROR; received from peer
func PCRJP(c *http.Client) model.Result {
	name := "Princess Connect Re - Dive Japan"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api-priconne-redive.cygames.jp/"
	headers := map[string]string{
		"User-Agent": model.UA_Dalvik,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
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
		Err: fmt.Errorf("get api-priconne-redive.cygames.jp failed with code: %d", resp.StatusCode)}
}
