package uk

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
// disco-api.discoveryplus.co.uk 仅 ipv4 且 get 请求
func DiscoveryPlus(c *http.Client) model.Result {
	name := "Discovery+ UK"
	hostname := "discoveryplus.co.uk"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://disco-api.discoveryplus.co.uk/token?realm=questuk&deviceId=61ee588b07c4df08c02861ecc1366a592c4ad02d08e8228ecfee67501d98bf47&shortlived=true"
	client1 := utils.Req(c)
	resp1, err1 := client1.R().Get(url1)
	if err1 != nil {
		return utils.HandleNetworkError(c, hostname, err1, name)
	}
	defer resp1.Body.Close()
	b1, err1 := io.ReadAll(resp1.Body)
	if err1 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	//body1 := string(b1)
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
	headers := map[string]string{
		"Cookie": cookies,
	}
	url2 := "https://disco-api.discoveryplus.co.uk/users/me"
	client2 := utils.Req(c)
	client2 = utils.SetReqHeaders(client2, headers)
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
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		if exit {
			if loc == "gb" {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc, Info: "origin", UnlockType: unlockType}
			} else {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc, Info: "global", UnlockType: unlockType}
			}
		}
		return model.Result{Name: name, Status: model.StatusNo, Region: loc}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get disco-api.discoveryplus.co.uk failed with code: %d", resp2.StatusCode)}
}
