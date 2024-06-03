package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// FXNOW
// fxnow.fxnetworks.com 仅 ipv4 且 get 请求
func FXNOW(request *gorequest.SuperAgent) model.Result {
	name := "FXNOW"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://fxnow.fxnetworks.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "is not accessible") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "FX Movies") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get fxnow.fxnetworks.com with code: %d", resp.StatusCode)}
}
