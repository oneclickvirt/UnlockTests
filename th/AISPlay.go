package th

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// AISPlay
// 49-231-37-237-rewriter.ais-vidnt.com 双栈 get 请求
func AISPlay(request *gorequest.SuperAgent) model.Result {
	name := "AIS Play"
	url := "https://49-231-37-237-rewriter.ais-vidnt.com/ais/play/origin/VOD/playlist/ais-yMzNH1-bGUxc/index.m3u8"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		if strings.Contains(body, "X-Geo-Protection-System-Status") {
			if strings.Contains(body, "ALLOW") {
				return model.Result{Name: name, Status: model.StatusYes}
			} else if strings.Contains(body, "BLOCK") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get 49-231-37-237-rewriter.ais-vidnt.com failed with code: %d", resp.StatusCode)}
}
