package us

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Philo
// content-us-east-2-fastly-b.www.philo.com 仅 ipv4 且 get 请求
func Philo(request *gorequest.SuperAgent) model.Result {
	name := "Philo"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://content-us-east-2-fastly-b.www.philo.com/geo"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Status  string `json:"status"`
		Country string `json:"country"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if resp.StatusCode == 403 || resp.StatusCode == 451 {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Status == "FAIL" {
		return model.Result{Name: name, Status: model.StatusNo, Region: strings.ToLower(res.Country)}
	} else if res.Status == "SUCCESS" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(res.Country)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get content-us-east-2-fastly-b.www.philo.com failed with code: %d", resp.StatusCode)}
}
