package hk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// ViuTV
// api.viu.now.com 双栈 且 post 请求
func ViuTV(c *http.Client) model.Result {
	name := "Viu.TV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.viu.now.com/p8/3/getLiveURL"
	resp, body, errs := utils.PostJson(c, url,
		"{\"callerReferenceNo\":\"20210726112323\",\"contentId\":\"099\",\"contentType\":\"Channel\",\"channelno\":\"099\",\"mode\":\"prod\",\"deviceId\":\"29b3cb117a635d5b56\",\"deviceType\":\"ANDROID_WEB\"}",
		map[string]string{"User-Agent": model.UA_Browser})
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		ResponseCode string
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.ResponseCode == "SUCCESS" {
		return model.Result{Status: model.StatusYes}
	} else if res.ResponseCode == "GEO_CHECK_FAIL" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.viu.now.com failed with code: %d", resp.StatusCode)}
}
