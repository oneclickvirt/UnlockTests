package fr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// FranceTV
// ftven.fr 纯 IPV4 get 请求
func FranceTV(c *http.Client) model.Result {
	name := "France TV"
	hostname := "ftven.fr"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://geo-info.ftven.fr/ws/edgescape.json"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	// body := string(b)
	// fmt.Println(body)
	var res struct {
		Response struct {
			GeoInfo struct {
				CountryCode string `json:"country_code"`
			} `json:"geo_info"`
		} `json:"reponse"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Response.GeoInfo.CountryCode == "FR" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// return model.Result{Name: name, Status: model.StatusUnexpected,
	// 	Err: fmt.Errorf("get canalplus.com failed with code: %d", resp.StatusCode)}
}
