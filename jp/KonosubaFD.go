package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// KonosubaFD
// api.konosubafd.jp 仅 ipv4 且 post 请求
func KonosubaFD(c *http.Client) model.Result {
	name := "Konosuba Fantastic Days"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.konosubafd.jp/api/masterlist"
	headers := map[string]string{
		"User-Agent": "pj0007/212 CFNetwork/1240.0.4 Darwin/20.6.0",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Post(url).End()
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
