package ca

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Crave
// capi.9c9media.com 仅 ipv4 且 get 请求
func Crave(request *gorequest.SuperAgent) model.Result {
	name := "Crave"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://capi.9c9media.com/destinations/se_atexace/platforms/desktop/bond/contents/2205173/contentpackages/4279732/manifest.mpd"
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Geo Constraint Restrictions") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
