package fr

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// CanalPlus
// canalplus.com 双栈 get 请求
func CanalPlus(c *http.Client) model.Result {
	name := "Canal+"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://boutique-tunnel.canalplus.com/"
	request := utils.Gorequest(c)
	resp, body, errs := request.Get(url).End()
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
