package ca

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Crave
// capi.9c9media.com 仅 ipv4 且 get 请求
func Crave(c *http.Client) model.Result {
	name := "Crave"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://capi.9c9media.com/destinations/se_atexace/platforms/desktop/bond/contents/2205173/contentpackages/4279732/manifest.mpd"
	request := utils.Gorequest(c)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Geo Constraint Restrictions") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(body, "video.9c9media.com") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get capi.9c9media.com with code: %d", resp.StatusCode)}
}
