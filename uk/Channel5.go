package uk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"time"
)

// Channel5
// cassie.channel5.com 仅 ipv4 且 get 请求
func Channel5(request *gorequest.SuperAgent) model.Result {
	name := "Channel 5"
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("https://cassie.channel5.com/api/v2/live_media/my5desktopng/C5.json?timestamp=%d&auth=0_rZDiY0hp_TNcDyk2uD-Kl40HqDbXs7hOawxyqPnbI", timestamp)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.code == "3000" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.code == "4003" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get cassie.channel5.com failed with code: %d", resp.StatusCode)}
}
