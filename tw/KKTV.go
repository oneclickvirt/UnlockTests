package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// KKTV
// api.kktv.me 仅 ipv4 且 get 请求
func KKTV(request *gorequest.SuperAgent) model.Result {
	name := "KKTV"
	url := "https://api.kktv.me/v3/ipcheck"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Data struct {
			Country   string `json:"country"`
			IsAllowed bool   `json:"is_allowed"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "\"is_allowed\":false") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Data.Country == "TW" && res.Data.IsAllowed {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if !res.Data.IsAllowed {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.kktv.me failed with head: %d", resp.StatusCode)}
}
