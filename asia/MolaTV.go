package asia

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// MolaTV
// mola.tv 仅 ipv4 且 get 请求
func MolaTV(request *gorequest.SuperAgent) model.Result {
	name := "Mola TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://mola.tv/api/v2/videos/geoguard/check/vd30491025"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var res struct {
		Data struct {
			Type       string `json:"type"`
			Id         string `json:"id"`
			Attributes struct {
				IsAllowed bool `json:"isAllowed"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "\"isAllowed\":false") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	if !res.Data.Attributes.IsAllowed {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if res.Data.Attributes.IsAllowed {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mola.tv failed with code: %d", resp.StatusCode)}
}
