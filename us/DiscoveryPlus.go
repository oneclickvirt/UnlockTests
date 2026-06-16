package us

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// DiscoveryPlus
// discoveryplus.com 双栈 且 post 请求
func DiscoveryPlus(c *http.Client) model.Result {
	name := "Discovery+"
	hostname := "discoveryplus.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://us1-prod-direct.discoveryplus.com/token?" +
		"deviceId=d1a4a5d25212400d1e6985984604d740&realm=go&shortlived=true"
	client1 := utils.Req(c)
	resp1, err1 := client1.R().Get(url1)
	if err1 != nil {
		return utils.HandleNetworkError(c, hostname, err1, name)
	}
	defer resp1.Body.Close()
	b1, err1 := io.ReadAll(resp1.Body)
	if err1 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err1}
	}
	var res struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b1, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusNo, Err: err}
	}
	cookies := "st=" + res.Data.Attributes.Token
	url2 := "https://us1-prod-direct.discoveryplus.com/users/me"
	headers2 := map[string]string{
		"Cookie": cookies,
	}
	client2 := utils.Req(c)
	client2 = utils.SetReqHeaders(client2, headers2)
	resp2, err2 := client2.R().Get(url2)
	if err2 != nil {
		return utils.HandleNetworkError(c, hostname, err2, name)
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	var res2 struct {
		Data struct {
			Attributes struct {
				CurrentLocationTerritory string `json:"currentLocationTerritory"`
			} `json:"attributes"`
		} `json:"data"`
	}
	//fmt.Println(string(b2))
	if err = json.Unmarshal(b2, &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if res2.Data.Attributes.CurrentLocationTerritory != "" {
		loc := strings.ToLower(res2.Data.Attributes.CurrentLocationTerritory)
		exit := utils.GetRegion(loc, model.DiscoveryPlusSupportCountry)
		if exit {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			if loc == "us" {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc, Info: "origin", UnlockType: unlockType}
			} else {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc, Info: "global", UnlockType: unlockType}
			}
		}
		return model.Result{Name: name, Status: model.StatusNo, Region: loc}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get us1-prod-direct.discoveryplus.com failed with code: %d", resp2.StatusCode)}
}
