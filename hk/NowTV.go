package hk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// NowTV
// webtvapi.now.com 仅 ipv4 且 post 请求
func NowTV(c *http.Client) model.Result {
	name := "Now TV"
	hostname := "now.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://webtvapi.now.com/10/7/getLiveURL"
	data1 := `{"contentId":"332","contentType":"Channel","deviceType":"IOS_PHONE","deviceId":"8269809F-7702-45CE-9378-D7157A2E6819","callerReferenceNo":"20140702122500","mode":"prod"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, body, err := utils.PostJson(c, url1, data1, headers)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	var res struct {
		ResponseCode string `json:"responseCode"` // 主要字段
	}
	// fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	switch res.ResponseCode {
	case "SUCCESS", "NOT_LOGIN", "ASSET_MISSING":
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	case "GEO_CHECK_FAIL":
		return model.Result{Name: name, Status: model.StatusNo}
	default:
		return model.Result{
			Name:   name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("webtvapi.now.com get unexpected responseCode: %s", res.ResponseCode),
		}
	}
}
