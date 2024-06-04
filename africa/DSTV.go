package africa

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/ecs/mediatest/utils"
	"net/http"
)

// DSTV
// authentication.dstv.com 仅 ipv4 且 get 请求
func DSTV(c *http.Client) model.Result {
	name := "DSTV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://authentication.dstv.com/favicon.ico"
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
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get authentication.dstv.com failed with code: %d", resp.StatusCode)}
}

