package nz

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ThreeNow
// bravo-livestream.fullscreen.nz 仅 ipv4 且 get 请求
func ThreeNow(request *gorequest.SuperAgent) model.Result {
	name := "ThreeNow"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://bravo-livestream.fullscreen.nz/index.m3u8"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Access Denied") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get bravo-livestream.fullscreen.nz failed with code: %d", resp.StatusCode)}
}
