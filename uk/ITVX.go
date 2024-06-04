package uk

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// ITVX
// simulcast.itv.com 仅 ipv4 且 get 请求
func ITVX(c *http.Client) model.Result {
	name := "ITV Hub"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://simulcast.itv.com/playlist/itvonline/ITV"
	request := utils.Gorequest(c)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || strings.Contains(body, "Outside Of Allowed Geographic Region") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(body, "Playlist") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get simulcast.itv.com failed with code: %d", resp.StatusCode)}
}
