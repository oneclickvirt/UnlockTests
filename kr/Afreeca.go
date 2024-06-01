package kr

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Afreeca
// vod.afreecatv.com 仅 ipv4 且 get 请求
func Afreeca(request *gorequest.SuperAgent) model.Result {
	name := "Afreeca TV"
	url := "https://vod.afreecatv.com/player/97464151"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if !strings.Contains(body, "document.location.href='https://vod.afreecatv.com'") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
