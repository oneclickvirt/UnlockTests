package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// AcornTV
// acorn.tv 仅 ipv4 且 get 请求
func AcornTV(request *gorequest.SuperAgent) model.Result {
	name := "Acorn TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://acorn.tv/"
	resp, body, errs := request.Get(url).Retry(2, 10).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Not yet available in your country") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get acorn.tv failed with code: %d", resp.StatusCode)}
}
