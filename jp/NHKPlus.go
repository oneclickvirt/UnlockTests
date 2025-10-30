package jp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// NHKPlus
// location-plus.nhk.jp 双栈 get 请求
func NHKPlus(c *http.Client) model.Result {
	name := "NHK+"
	hostname := "nhk.jp"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://location-plus.nhk.jp/geoip/area.json"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	var res struct {
		CountryCode string `json:"country_code"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.CountryCode == "JP" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo, Err: err}
	// return model.Result{Name: name, Status: model.StatusUnexpected,
	// 	Err: fmt.Errorf("get location-plus.nhk.jp failed with code: %d", resp.StatusCode)}
}
