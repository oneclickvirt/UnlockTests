package it

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// RaiPlay
// mediapolisvod.rai.it 仅 ipv4 且 get 请求
func RaiPlay(request *gorequest.SuperAgent) model.Result {
	name := "Rai Play"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://mediapolisvod.rai.it/relinker/relinkerServlet.htm?cont=VxXwi7UcqjApssSlashbjsAghviAeeqqEEqualeeqqEEqual&output=64"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "video_no_available") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		// TODO 识别地区
		// <geoprotection>
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mediapolisvod.rai.it failed with code: %d", resp.StatusCode)}
}
