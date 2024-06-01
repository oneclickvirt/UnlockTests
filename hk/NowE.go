package hk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NowE
// webtvapi.nowe.com 仅 ipv4 且 post 请求
func NowE(request *gorequest.SuperAgent) model.Result {
	name := "Now E"
	url1 := "https://webtvapi.nowe.com/16/1/getVodURL"
	data1 := `{"contentId":"202403181904703","contentType":"Vod","pin":"","deviceName":"Browser","deviceId":"w-663bcc51-913c-913c-913c-913c913c","deviceType":"WEB","secureCookie":null,"callerReferenceNo":"W17151951620081575","profileId":null,"mupId":null}`
	resp, body, errs := request.Post(url1).
		Send(data1).
		Set("Content-Type", "application/json").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		ResponseCode string `json:"responseCode"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.ResponseCode == "GEO_CHECK_FAIL" {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if res.ResponseCode == "SUCCESS" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{
		Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("webtvapi.nowe.com get responseCode: %s", res.ResponseCode),
	}
}
