package transnation

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// SonyLiv
// www.sonyliv.com 双栈 且 get 请求
func SonyLiv(c *http.Client) model.Result {
	name := "SonyLiv"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.sonyliv.com/"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp1, body1, errs1 := request.Get(url).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	if strings.Contains(body1, "geolocation_notsupported") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Unavailable"}
	}
	jwtToken := utils.ReParse(body1, `resultObj:"([^"]+)`)

	headers2 := map[string]string{
		"accept":         "application/json, text/plain, */*",
		"referer":        "https://www.sonyliv.com/",
		"device_id":      "25a417c3b5f246a393fadb022adc82d5-1715309762699",
		"app_version":    "3.5.59",
		"security_token": jwtToken,
	}
	url2 := "https://apiv2.sonyliv.com/AGL/1.4/A/ENG/WEB/ALL/USER/ULD"
	request2 := utils.Gorequest(c)
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	resp2, body2, errs2 := request2.Get(url2).End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	var res1 struct {
		ResultObj struct {
			CountryCode string `json:"country_code"`
		} `json:"resultObj"`
	}
	if err := json.Unmarshal([]byte(body2), &res1); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	region := res1.ResultObj.CountryCode
	if region == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("can not found region")}
	}

	headers3 := map[string]string{
		"upgrade-insecure-requests": "1",
		"accept":                    "application/json, text/plain, */*",
		"origin":                    "https://www.sonyliv.com",
		"referer":                   "https://www.sonyliv.com/",
		"device_id":                 "25a417c3b5f246a393fadb022adc82d5-1715309762699",
		"security_token":            jwtToken,
	}
	url3 := "https://apiv2.sonyliv.com/AGL/3.8/A/ENG/WEB/" + region + "/ALL/CONTENT/VIDEOURL/VOD/1000273613/prefetch"
	request3 := utils.Gorequest(c)
	request3 = utils.SetGoRequestHeaders(request3, headers3)
	resp3, body3, errs3 := request3.Get(url3).End()
	if len(errs3) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs3[0]}
	}
	defer resp3.Body.Close()
	var res2 struct {
		ResultCode string `json:"resultCode"`
	}
	if err := json.Unmarshal([]byte(body3), &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res2.ResultCode == "OK" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	}
	if res2.ResultCode == "KO" {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Proxy Detected", Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get apiv2.sonyliv.com failed with code: %d", resp3.StatusCode)}
}
