package uk

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// BBCiPlayer
// open.live.bbc.co.uk 仅 ipv4 且 get 请求
func BBCiPlayer(c *http.Client) model.Result {
	name := "BBC iPLAYER"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://open.live.bbc.co.uk/mediaselector/6/select/version/2.0/mediaset/pc/vpid/bbc_one_london/format/json/jsfunc/JS_callbacks0"
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
	if resp.StatusCode == 200 {
		if strings.Contains(body, "geolocation") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if strings.Contains(body, "vs-hls-push-uk") {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	} else if resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get open.live.bbc.co.uk failed with code: %d", resp.StatusCode)}
}
