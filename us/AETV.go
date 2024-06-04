package us

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// AETV
// ccpa-service.sp-prod.net 仅 ipv4 且 post 请求
func AETV(c *http.Client) model.Result {
	name := "A&E TV"
	if c == nil {
		return model.Result{Name: name}
	}

	url1 := "https://link.theplatform.com/s/xc6n8B/UR27JDU0bu2s/"
	headers1 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request1 := utils.Gorequest(c)
	request1 = utils.SetGoRequestHeaders(request1, headers1)
	resp1, body1, errs1 := request1.Post(url1).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	if strings.Contains(body1, "GeoLocationBlocked") {
		return model.Result{Name: name, Status: model.StatusNo}
	}

	url2 := "https://play.aetv.com/"
	headers2 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request2 := utils.Gorequest(c)
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	resp2, body2, errs2 := request2.Post(url2).End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	if body2 != "" {
		tp := utils.ReParse(body2, "AETN-Country-Code=\\K[A-Z]+")
		if tp != "" {
			region := strings.ToLower(tp)
			if region == "ca" || region == "us" {
				return model.Result{Name: name, Status: model.StatusYes, Region: region}
			} else {
				return model.Result{Name: name, Status: model.StatusNo}
			}
		}
	}

	url3 := "https://ccpa-service.sp-prod.net/ccpa/consent/10265/display-dns"
	headers3 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request3 := utils.Gorequest(c)
	request3 = utils.SetGoRequestHeaders(request3, headers3)
	resp3, body3, errs3 := request3.Post(url3).End()
	if len(errs3) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs3[0]}
	}
	defer resp3.Body.Close()
	//fmt.Println(body)
	var res struct {
		CcpaApplies bool `json:"ccpaApplies"`
	}
	if err := json.Unmarshal([]byte(body3), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.CcpaApplies == true {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if res.CcpaApplies == false {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ccpa-service.sp-prod.net failed with code: %d", resp3.StatusCode)}
}
