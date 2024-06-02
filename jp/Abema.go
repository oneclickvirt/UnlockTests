package jp

import (
	"encoding/json"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Abema
// api.abema.io 仅 ipv4 且 get 请求
func Abema(request *gorequest.SuperAgent) model.Result {
	name := "Abema.TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api.abema.io/v1/ip/check?device=android"
	request = request.Set("User-Agent", model.UA_Dalvik)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var abemaRes struct {
		IsoCountryCode string `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &abemaRes); err != nil {
		if strings.Contains(body, "blocked_location") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if abemaRes.IsoCountryCode == "JP" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if abemaRes.IsoCountryCode == "blocked_location" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes + " (Oversea Only)"}
}
