package eu

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// SetantaSports
// dce-frontoffice.imggaming.com 仅 ipv4 且 get 请求
func SetantaSports(request *gorequest.SuperAgent) model.Result {
	name := "Setanta Sports"
	url := "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt"
	resp, body, errs := request.Get(url).
		Set("User-Agent", model.UA_Browser).
		Set("Realm", "dce.adjara").
		Set("x-api-key", "857a1e5d-e35e-4fdf-805b-a87b6f8364bf").
		Set("Accept-Language", "en-US").
		Retry(2, 5).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var consentResponse struct {
		OutsideAllowedTerritories bool `json:"outsideAllowedTerritories"`
	}
	if err := json.Unmarshal([]byte(body), &consentResponse); err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	if consentResponse.OutsideAllowedTerritories {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if !consentResponse.OutsideAllowedTerritories {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{
		Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get dce-frontoffice.imggaming.com failed with code: %d", resp.StatusCode)}
}
