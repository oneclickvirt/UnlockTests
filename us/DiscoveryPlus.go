package us

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
	"io"
	"strings"
)

// DiscoveryPlus
// discoveryplus.com 双栈 且 post 请求
func DiscoveryPlus(request *gorequest.SuperAgent) model.Result {
	name := "Discovery+"
	url1 := "https://us1-prod-direct.discoveryplus.com/token?" +
		"deviceId=d1a4a5d25212400d1e6985984604d740&realm=go&shortlived=true"
	resp1, _, errs1 := request.Get(url1).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	b1, err1 := io.ReadAll(resp1.Body)
	if err1 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err1}
	}
	var res struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b1, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusNo, Err: err}
	}
	cookies := "_gcl_au=1.1.858579665.1632206782; _rdt_uuid=1632206782474.6a9ad4f2-8ef7-4a49-9d60-e071bce45e88; " +
		"_scid=d154b864-8b7e-4f46-90e0-8b56cff67d05; " +
		"_pin_unauth=dWlkPU1qWTRNR1ZoTlRBdE1tSXdNaTAwTW1Nd0xUbGxORFV0WWpZMU0yVXdPV1l6WldFeQ; " +
		"_sctr=1|1632153600000; aam_fw=aam%3D9354365%3Baam%3D9040990; " +
		fmt.Sprintf("aam_uuid=24382050115125439381416006538140778858; st=%s; ", res.Data.Attributes.Token) +
		"gi_ls=0; _uetvid=a25161a01aa711ec92d47775379d5e4d; " +
		"AMCV_BC501253513148ED0A490D45%40AdobeOrg=-1124106680%7CMCIDTS%7C18894%7CMCMID%7C24223296309793" +
		"747161435877577673078228%7CMCAAMLH-1633011393%7C9%7CMCAAMB-1633011393%7CRKhpRz8krg2tLO6pguXWp5o" +
		"lkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1632413793s%7CNONE%7CvVersion%7C5.2.0; " +
		"ass=19ef15da-95d6-4b1d-8fa2-e9e099c9cc38.1632408400.1632406594"
	resp2, _, errs2 := request.
		Get("https://us1-prod-direct.discoveryplus.com/users/me").
		Set("Cookie", cookies).
		End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	var res2 struct {
		Data struct {
			Attributes struct {
				CurrentLocationTerritory string `json:"currentLocationTerritory"`
			} `json:"attributes"`
		} `json:"data"`
	}
	//fmt.Println(string(b2))
	if err := json.Unmarshal(b2, &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if res2.Data.Attributes.CurrentLocationTerritory != "" {
		loc := strings.ToLower(res2.Data.Attributes.CurrentLocationTerritory)
		exit := utils.GetRegion(loc, model.DiscoveryPlusSupportCountry)
		if exit {
			if loc == "us" {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc + " (origin)"}
			} else {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc + " (global)"}
			}
		}
		return model.Result{Name: name, Status: model.StatusNo, Region: loc}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get us1-prod-direct.discoveryplus.com failed with code: %d", resp2.StatusCode)}
}
