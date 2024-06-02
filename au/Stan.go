package au

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// Stan
// api.stan.com.au 仅 ipv4 且 post 请求
func Stan(request *gorequest.SuperAgent) model.Result {
	name := "Stan"
	resp, body, errs := utils.PostJson(request, "https://api.stan.com.au/login/v1/sessions/web/account", "{}")
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(string(body), "Access Denied") || resp.StatusCode == 404 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Unavailable"}
	}
	if strings.Contains(string(body), "VPNDetected") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "VPN Detected"}
	}
	if resp.StatusCode == 400 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.stan.com.au failed with code: %d", resp.StatusCode)}
}
