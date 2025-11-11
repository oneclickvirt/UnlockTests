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
		return utils.HandleNetworkError(c, hostname, err, name)
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
		return utils.HandleNetworkError(c, hostname, err, name)
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
		return utils.HandleNetworkError(c, hostname, err, name)
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
		return utils.HandleNetworkError(c, hostname, err, name)
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
	if region != "" && strings.Contains(availableRegion, region) {
		// VPN检测
		vpnCheckResp, err := client.R().
			SetHeader("Content-Type", "application/x-www-form-urlencoded").
			SetBody("st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi0wOWQxOTg4Yy1mZmUzLTQxMDEtOWI5My0yNDU1ZTkyNGQ1YjYiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6YjYzOTgxZWQtNzA2MC00ZGYwLThkZGItZjA2YjFkNWRjZWVkIiwiaWF0IjoxNzQzODQwMzgwLCJleHAiOjIwNTkyMDAzODAsInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fYW1lciIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6IjQwYTgzZjNlLTY4OTktNDE3Mi1hMWY2LWJjZDVjN2ZkNjA4NSIsInZlcnNpb24iOiJ2MyIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiNWY3YzViZjQtYjc4Ny00NDRjLWJhYTYtMzU5MzgwYWFiM2RmIn0.f5HTgIV2v0nQQDp5LQG0xqLrxyACdvnMDiWO_viX_CUGqtc5ncSjp_LgM30QFkkMnINFhzKEGRpsZvb-o3Pj_Z39uRBr5LCeiCPR7ssV-_SXyRFVRRDEB2lpxyz7jmdD1SxvA06HnEwTbZQzlbZ7g9GXq02yNdEfHlqYEh_4WF88UbXfeieYTd4TH7kwN1RE50NfQUS6f0WmzpAbpiULyd87mpTeynchFNMMz-YHVzZ_-nDW6geihXc3tS0FKVSR8fdOSPQFzEYOLCfhInufiPahiXI-OKF89aShAqM-y4Hx_eukGnsq3mO5wa3unnqVr9Kzc61BIhHh1Hs2bqYiYg").
			Post("https://default.any-any.prd.api.max.com/any/playback/v1/playbackInfo")
		if err != nil {
			return utils.HandleNetworkError(c, hostname, err, name)
		}
		vpnCheckBody, err := io.ReadAll(vpnCheckResp.Body)
		if err != nil {
			return model.Result{Status: model.StatusNetworkErr, Err: err}
		}
		if strings.Contains(string(vpnCheckBody), "VPN") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{
			Name:       name,
			Status:     model.StatusYes,
			Region:     strings.ToLower(region),
			UnlockType: unlockType,
		}
	}
	return model.Result{
		Name:   name,
		Status: model.StatusNo,
		Region: strings.ToLower(region),
	}
}
