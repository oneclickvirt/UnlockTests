package sg

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// MeWatch
// cdn.mewatch.sg 仅 ipv4 且 get 请求
func MeWatch(request *gorequest.SuperAgent) model.Result {
	name := "MeWatch"
	url := "https://cdn.mewatch.sg/api/items/97098/videos?delivery=stream%2Cprogressive&ff=idp%2Cldp%2Crpt%2Ccd&lang=en&resolution=External&segments=all"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if strings.Contains(body, "You are accessing this item from a location that is not permitted by the license") ||
		res.Code == 8002 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || strings.Contains(body, "deliveryType") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get cdn.mewatch.sg failed with code: %d", resp.StatusCode)}
}
