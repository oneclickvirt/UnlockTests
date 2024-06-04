package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// FXNOW
// fxnow.fxnetworks.com 仅 ipv4 且 get 请求
func FXNOW(c *http.Client) model.Result {
	name := "FXNOW"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://fxnow.fxnetworks.com/"
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
	if strings.Contains(body, "is not accessible") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "FX Movies") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get fxnow.fxnetworks.com with code: %d", resp.StatusCode)}
}
