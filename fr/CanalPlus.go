package fr

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// CanalPlus
// canalplus.com 双栈 get 请求
func CanalPlus(request *gorequest.SuperAgent) model.Result {
	name := "Canal+"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://boutique-tunnel.canalplus.com/"
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "othercountry") ||
		strings.Contains(resp.Request.URL.String(), "other-country-blocking") ||
		resp.StatusCode == 302 || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get canalplus.com failed with code: %d", resp.StatusCode)}
}
