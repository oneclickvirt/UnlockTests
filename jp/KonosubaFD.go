package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// KonosubaFD
// api.konosubafd.jp 仅 ipv4 且 post 请求
func KonosubaFD(request *gorequest.SuperAgent) model.Result {
	name := "Konosuba Fantastic Days"
	url := "https://api.konosubafd.jp/api/masterlist"
	request = request.Set("User-Agent", "pj0007/212 CFNetwork/1240.0.4 Darwin/20.6.0")
	resp, _, errs := request.Post(url).Retry(1, 30).End()
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
		Err: fmt.Errorf("get api.konosubafd.jp failed with code: %d", resp.StatusCode)}
}
