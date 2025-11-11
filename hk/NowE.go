package hk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// NowE
// webtvapi.nowe.com 仅 ipv4 且 post 请求
func NowE(c *http.Client) model.Result {
	name := "Now E"
	hostname := "nowe.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://webtvapi.nowe.com/16/1/getVodURL"
	data1 := `{"contentId":"202310181863841","contentType":"Vod","pin":"","deviceName":"Browser","deviceId":"w-678913af-3998-3998-3998-39983998","deviceType":"WEB","secureCookie":null,"callerReferenceNo":"W17370372345461425","profileId":null,"mupId":null,"trackId":"738296446.226.1737037103860.2","sessionId":"c39f03e6-9e74-4d24-a82f-e0d0f328bb70"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, body, err := utils.PostJson(c, url1, data1, headers)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	var res struct {
		ResponseCode        string `json:"responseCode"`         // 主要字段
		OTTAPIResponseCode  string `json:"OTTAPI_ResponseCode"` // 备选字段
	}
	// fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.OTTAPIResponseCode == "SUCCESS" || res.ResponseCode == "NOT_LOGIN" || res.ResponseCode == "ASSET_MISSING" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	} else if res.ResponseCode == "GEO_CHECK_FAIL" {
		return model.Result{Name: name, Status: model.StatusNo}
	} else {
		return model.Result{
			Name:   name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("webtvapi.nowe.com get unexpected responseCode: %s", res.ResponseCode),
		}
	}
}