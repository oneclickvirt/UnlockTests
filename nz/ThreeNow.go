package nz

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// ThreeNow
// bravo-livestream.fullscreen.nz 仅 ipv4 且 get 请求
func ThreeNow(c *http.Client) model.Result {
	name := "ThreeNow"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://bravo-livestream.fullscreen.nz/index.m3u8"
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
	if strings.Contains(body, "Access Denied") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get bravo-livestream.fullscreen.nz failed with code: %d", resp.StatusCode)}
}
