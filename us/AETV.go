package us

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// AETV
// ccpa-service.sp-prod.net 仅 ipv4 且 post 请求
func AETV(request *gorequest.SuperAgent) model.Result {
	name := "A&E TV"
	url := "https://ccpa-service.sp-prod.net/ccpa/consent/10265/display-dns"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Post(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var res struct {
		CcpaApplies bool `json:"ccpaApplies"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.CcpaApplies == true {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if res.CcpaApplies == false {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ccpa-service.sp-prod.net failed with code: %d", resp.StatusCode)}
}
