package transnation

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// TVBAnywhere
// uapisfm.tvbanywhere.com.sg 仅 ipv4 且 get 请求
func TVBAnywhere(request *gorequest.SuperAgent) model.Result {
	name := "TVBAnywhere+"
	url := "https://uapisfm.tvbanywhere.com.sg/geoip/check/platform/android"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		AllowInThisCountry bool   `json:"allow_in_this_country"`
		Country            string `json:"country"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.AllowInThisCountry && res.Country != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(res.Country)}
	} else if !res.AllowInThisCountry && res.Country != "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get uapisfm.tvbanywhere.com.sg failed with code: %d", resp.StatusCode)}
}
