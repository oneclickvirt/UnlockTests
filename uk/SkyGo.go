package uk

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// SkyGo
// skyid.sky.com 仅 ipv4 且 get 请求
func SkyGo(c *http.Client) model.Result {
	name := "Sky Go"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://skyid.sky.com/authorise/skygo?response_type=token&client_id=sky&appearance=compact&redirect_uri=skygo://auth"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "You don't have permission to access") || resp.StatusCode == 403 || resp.StatusCode == 200 ||
		strings.Contains(body, "Access Denied") { // || resp.StatusCode == 451
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get skyid.sky.com failed with code: %d", resp.StatusCode)}
}
