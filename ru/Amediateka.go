package ru

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Amediateka
// www.amediateka.ru 仅 ipv4 且 get 请求
func Amediateka(request *gorequest.SuperAgent) model.Result {
	name := "Amediateka"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.amediateka.ru/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "VPN") || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 301 && resp.Header.Get("Location") == "https://www.amediateka.ru/unavailable/index.html" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.amediateka.ru failed with code: %d", resp.StatusCode)}
}
