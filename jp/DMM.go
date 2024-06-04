package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// DMM
// bitcoin.dmm.com 仅 ipv4 且 get 请求
func DMM(c *http.Client) model.Result {
	name := "DMM"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://bitcoin.dmm.com"
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
	if strings.Contains(body, "This page is not available in your area") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "暗号資産") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get bitcoin.dmm.com failed with code: %d", resp.StatusCode)}
}
