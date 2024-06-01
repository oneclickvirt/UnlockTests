package us

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Crunchyroll
// c.evidon.com 仅 ipv4 且 get 请求
func Crunchyroll(request *gorequest.SuperAgent) model.Result {
	name := "Crunchyroll"
	url := "https://c.evidon.com/geo/country.js"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 10).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "'code':'us'") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
