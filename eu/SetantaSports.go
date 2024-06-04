package eu

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// SetantaSports
// dce-frontoffice.imggaming.com 仅 ipv4 且 get 请求
func SetantaSports(c *http.Client) model.Result {
	name := "Setanta Sports"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt"
	headers := map[string]string{
		"User-Agent":      model.UA_Browser,
		"Realm":           "dce.adjara",
		"x-api-key":       "857a1e5d-e35e-4fdf-805b-a87b6f8364bf",
		"Accept-Language": "en-US",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
