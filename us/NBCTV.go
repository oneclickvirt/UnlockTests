package us

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// NBCTV
// geolocation.digitalsvc.apps.nbcuni.com 双栈 get 请求
func NBCTV(c *http.Client) model.Result {
	name := "NBC TV"
	hostname := "nbcuni.com"
	if c == nil {
		return model.Result{Name: name}
	}
	fakeUuid, _ := uuid.NewV4()
	url := "https://geolocation.digitalsvc.apps.nbcuni.com/geolocation/live/usa"
	client := utils.Req(c)
	headers := map[string]string{
		"accept-language":    "en-US,en;q=0.9",
		"app-session-id":     fakeUuid.String(),
		"authorization":      "NBC-Basic key=\"usa_live\", version=\"3.0\", type=\"cpc\"",
		"client":             "oneapp",
		"content-type":       "application/json",
		"origin":             "https://www.nbc.com",
		"referer":            "https://www.nbc.com/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
	}
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().SetBodyJsonString(`{"adobeMvpdId":null,"serviceZip":null,"device":"web"}`).Post(url)
	if err == nil {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			var res struct {
				Restricted  bool `json:"restricted"`
				RequestInfo struct {
					CountryCode string `json:"countryCode"`
				} `json:"requestInfo"`
				RestrictionDetails struct {
					Code string `json:"code"`
				} `json:"restrictionDetails"`
			}
			if jsonErr := json.Unmarshal(b, &res); jsonErr == nil {
				if !res.Restricted {
					result1, result2, result3 := utils.CheckDNS(hostname)
					unlockType := utils.GetUnlockType(result1, result2, result3)
					return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
				} else if res.RestrictionDetails.Code == "321" {
					return model.Result{Name: name, Status: model.StatusNo}
				}
			} else {
				body := string(b)
				if body != "" && body != "{}" {
					if body == `{"restricted":false}` || (len(body) > 20 && body[1:12] == `"restricted"` && body[13:18] == "false") {
						result1, result2, result3 := utils.CheckDNS(hostname)
						unlockType := utils.GetUnlockType(result1, result2, result3)
						return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
					} else if body == `{"restricted":true}` || (len(body) > 20 && body[1:12] == `"restricted"` && body[13:17] == "true") {
						return model.Result{Name: name, Status: model.StatusNo}
					}
				}
			}
		}
	}
	url = "https://geolocation.onetrust.com/cookieconsentpub/v1/geo/location/dnsfeed"
	client = utils.Req(c)
	resp, err = client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	var geoRes struct {
		Country string `json:"country"`
	}
	if jsonErr := json.Unmarshal(body, &geoRes); jsonErr == nil && geoRes.Country == "US" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("geolocation.digitalsvc.apps.nbcuni.com failed with code: %d", resp.StatusCode)}
}
