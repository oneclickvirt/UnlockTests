package uk

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BBCiPlayer
// open.live.bbc.co.uk 仅 ipv4 且 get 请求
func BBCiPlayer(request *gorequest.SuperAgent) model.Result {
	name := "BBC iPLAYER"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://open.live.bbc.co.uk/mediaselector/6/select/version/2.0/mediaset/pc/vpid/bbc_one_london/format/json/jsfunc/JS_callbacks0"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		if strings.Contains(body, "geolocation") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	} else if resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get open.live.bbc.co.uk failed with code: %d", resp.StatusCode)}
}
