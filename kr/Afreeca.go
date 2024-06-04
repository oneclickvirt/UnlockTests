package kr

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Afreeca
// vod.afreecatv.com 仅 ipv4 且 get 请求
func Afreeca(c *http.Client) model.Result {
	name := "Afreeca TV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://vod.afreecatv.com/player/97464151"
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
	if !strings.Contains(body, "document.location.href='https://vod.afreecatv.com'") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
