package us

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// AETV
// ccpa-service.sp-prod.net 仅 ipv4 且 post 请求
func AETV(c *http.Client) model.Result {
	name := "A&E TV"
	hostname := "aetv.com"
	if c == nil {
		return model.Result{Name: name}
	}

	url0 := "https://geo.privacymanager.io/"
	client0 := utils.Req(c)
	resp0, err0 := client0.R().Get(url0)
	if err0 == nil {
		defer resp0.Body.Close()
		b0, err0 := io.ReadAll(resp0.Body)
		if err0 == nil {
			var geoRes struct {
				Country string `json:"country"`
			}
			if err := json.Unmarshal(b0, &geoRes); err == nil {
				if geoRes.Country == "US" || geoRes.Country == "CA" {
					result1, result2, result3 := utils.CheckDNS(hostname)
					unlockType := utils.GetUnlockType(result1, result2, result3)
					return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(geoRes.Country), UnlockType: unlockType}
				}
			}
		}
	}

	url1 := "https://link.theplatform.com/s/xc6n8B/UR27JDU0bu2s/"
	client1 := utils.Req(c)
	resp1, err1 := client1.R().Post(url1)
	if err1 == nil {
		defer resp1.Body.Close()
		b1, err1 := io.ReadAll(resp1.Body)
		if err1 == nil {
			body1 := string(b1)
			if strings.Contains(body1, "GeoLocationBlocked") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
		}
	}

	url2 := "https://play.aetv.com/"
	client2 := utils.Req(c)
	resp2, err2 := client2.R().Post(url2)
	if err2 == nil {
		defer resp2.Body.Close()
		b2, err2 := io.ReadAll(resp2.Body)
		if err2 == nil {
			body2 := string(b2)
			if body2 != "" {
				tp := utils.ReParse(body2, `AETN-Country-Code=([A-Z]+)`)
				if tp != "" {
					region := strings.ToLower(tp)
					if region == "ca" || region == "us" {
						result1, result2, result3 := utils.CheckDNS(hostname)
						unlockType := utils.GetUnlockType(result1, result2, result3)
						return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
					} else {
						return model.Result{Name: name, Status: model.StatusNo}
					}
				}
			}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ccpa-service.sp-prod.net failed with code: %d", resp2.StatusCode)}
}
