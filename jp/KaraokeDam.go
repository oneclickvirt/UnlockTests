package jp

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// KaraokeDam
// cds1.clubdam.com 仅 ipv4 且 get 请求
func KaraokeDam(request *gorequest.SuperAgent) model.Result {
	name := "Karaoke@DAM"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "http://cds1.clubdam.com/vhls-cds1/site/xbox/sample_1.mp4.m3u8"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get Karaoke@Dam failed with code: %d", resp.StatusCode)}
}
