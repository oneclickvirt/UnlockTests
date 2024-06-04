package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// J_COM_ON_DEMAND
// linkvod.myjcom.jp 仅 ipv4 且 get 请求
func J_COM_ON_DEMAND(c *http.Client) model.Result {
	name := "J:com On Demand"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://linkvod.myjcom.jp/auth/login"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 404 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 502 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get linkvod.myjcom.jp failed with code: %d", resp.StatusCode)}
}
