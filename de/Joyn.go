package de

import (
	"encoding/json"
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// Joyn
// auth.joyn.de 仅 ipv4 且 post 请求
func Joyn(request *gorequest.SuperAgent) model.Result {
	name := "Joyn"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://auth.joyn.de/auth/anonymous"
	payload := `{"client_id":"b74b9f27-a994-4c45-b7eb-5b81b1c856e7","client_name":"web","anon_device_id":"b74b9f27-a994-4c45-b7eb-5b81b1c856e7"}`
	resp, body, errs := utils.PostJson(request, url, payload)
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}

	url2 := "https://api.joyn.de/content/entitlement-token"
	headers := map[string]string{
		"authorization": "Bearer " + res.AccessToken,
		"x-api-key":     "36lp1t4wto5uu2i2nk57ywy9on1ns5yg",
	}
	payload2 := `{"content_id":"daserste-de-hd","content_type":"LIVE"}`
	resp2, body2, errs2 := utils.PostJson(request, url2, payload2, headers)
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()

	var res2a []struct {
		Code string `json:"code"`
	}
	var res2b struct {
		Token string `json:"entitlement_token"`
	}
	if err := json.Unmarshal(body2, &res2a); err != nil {
		if err := json.Unmarshal(body2, &res2b); err != nil {
			return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
		}
		if res2b.Token != "" {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	}
	if len(res2a) > 0 && res2a[0].Code == "ENT_AssetNotAvailableInCountry" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get joyn.de with code: %d %d", resp.StatusCode, resp2.StatusCode)}
}
