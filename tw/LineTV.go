package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// LineTV
// www.linetv.tw 仅 ipv4 且 get 请求
func LineTV(c *http.Client) model.Result {
	name := "LineTV.TW"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.linetv.tw/api/part/11829/eps/1/part?chocomemberId="
	request := utils.Gorequest(c)
	resp, body, errs := request.Get(url).End()
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
