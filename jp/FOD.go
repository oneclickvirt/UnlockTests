package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// FOD
// geocontrol1.stream.ne.jp 仅 ipv4 且 get 请求
func FOD(request *gorequest.SuperAgent) model.Result {
	name := "FOD(Fuji TV)"
	url := "https://geocontrol1.stream.ne.jp/fod-geo/check.xml?time=1624504256"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "false") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || strings.Contains(body, "true") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get geocontrol1.stream.ne.jp failed with code: %d", resp.StatusCode)}
}
