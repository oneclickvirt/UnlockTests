package transnation

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// TVBAnywhere
// uapisfm.tvbanywhere.com.sg 仅 ipv4 且 get 请求
func TVBAnywhere(c *http.Client) model.Result {
	name := "TVBAnywhere+"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://uapisfm.tvbanywhere.com.sg/geoip/check/platform/android"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
