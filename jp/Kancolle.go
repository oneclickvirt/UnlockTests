package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Kancolle
// Kancolle 仅 ipv4 且 get 请求
func Kancolle(request *gorequest.SuperAgent) model.Result {
	name := "Kancolle Japan"
	url := "http://203.104.209.7/kcscontents/news/"
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
		Err: fmt.Errorf("get Kancolle failed with code: %d", resp.StatusCode)}
}
