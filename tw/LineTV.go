package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// LineTV
// www.linetv.tw 仅 ipv4 且 get 请求
func LineTV(request *gorequest.SuperAgent) model.Result {
	name := "LineTV.TW"
	url := "https://www.linetv.tw/api/part/11829/eps/1/part?chocomemberId="
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		CountryCode int `json:"countryCode"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.CountryCode == 228 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if resp.StatusCode == 400 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.linetv.tw failed with code: %d", resp.StatusCode)}
}
