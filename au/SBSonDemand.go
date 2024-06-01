package au

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// SBSonDemand
// www.sbs.com.au 仅 ipv4 且 get 请求
func SBSonDemand(request *gorequest.SuperAgent) model.Result {
	name := "SBS on Demand"
	url := "https://www.sbs.com.au/api/v3/network?context=odwebsite"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Get struct {
			Response struct {
				CountryCode string `json:"country_code"`
			} `json:"response"`
		} `json:"get"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Get.Response.CountryCode == "AU" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
