package us

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

type hboTokenResponse struct {
	Data struct {
		Attributes struct {
			Token string `json:"token"`
		} `json:"attributes"`
	} `json:"data"`
}

type hboSessionResponse struct {
	Routing struct {
		Domain     string `json:"domain"`
		Tenant     string `json:"tenant"`
		Env        string `json:"env"`
		HomeMarket string `json:"homeMarket"`
	} `json:"routing"`
}

type hboUserResponse struct {
	Data struct {
		Attributes struct {
			CurrentLocationTerritory string `json:"currentLocationTerritory"`
		} `json:"attributes"`
	} `json:"data"`
}

func HBOMax(c *http.Client) model.Result {
	name := "HBO Max"
	hostname := "max.com"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	// 获取Token
	tokenResp, err := client.R().
		SetHeaders(map[string]string{
			"x-device-info":  "beam/5.0.0 (desktop/desktop; Windows/10; afbb5daa-c327-461d-9460-d8e4b3ee4a1f/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)",
			"x-disco-client": "WEB:10:beam:5.2.1",
		}).
		Get("https://default.any-any.prd.api.max.com/token?realm=bolt&deviceId=afbb5daa-c327-461d-9460-d8e4b3ee4a1f")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	var tokenData hboTokenResponse
	if err := json.Unmarshal(tokenBody, &tokenData); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	token := tokenData.Data.Attributes.Token
	if token == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}

	// 获取Session信息
	sessionResp, err := client.R().
		SetHeader("Cookie", "st="+token).
		SetHeader("Content-Type", "application/json").
		Post("https://default.any-any.prd.api.max.com/session-context/headwaiter/v1/bootstrap")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	sessionBody, err := io.ReadAll(sessionResp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	var sessionData hboSessionResponse
	if err := json.Unmarshal(sessionBody, &sessionData); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	// 构建用户信息URL
	userURL := fmt.Sprintf("https://default.%s-%s.%s.%s/users/me",
		sessionData.Routing.Tenant,
		sessionData.Routing.HomeMarket,
		sessionData.Routing.Env,
		sessionData.Routing.Domain,
	)
	// 获取用户信息
	userResp, err := client.R().
		SetHeader("Cookie", "st="+token).
		Get(userURL)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	var userData hboUserResponse
	if err := json.Unmarshal(userBody, &userData); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	region := userData.Data.Attributes.CurrentLocationTerritory
	// 验证区域可用性
	checkResp, err := client.R().Get("https://www.max.com/")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	checkBody, err := io.ReadAll(checkResp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	body := string(checkBody)
	availableRegion := strings.ToUpper(strings.Join(
		strings.Fields(
			strings.Join(
				strings.Split(body, "\"url\":\"/")[1:], " ")),
		" "))
	status := model.StatusNo
	if region != "" && strings.Contains(availableRegion, region) {
		status = model.StatusYes
		// DNS检测
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{
			Name:       name,
			Status:     status,
			Region:     strings.ToLower(region),
			UnlockType: unlockType,
		}
	}
	return model.Result{
		Name:   name,
		Status: status,
		Region: strings.ToLower(region),
	}
}
