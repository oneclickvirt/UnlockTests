package transnation

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"regexp"
	"strings"
)

func parseSonyLivToken(body string) string {
	re := regexp.MustCompile(`resultObj:"([^"]+)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// SonyLiv
// www.sonyliv.com 双栈 且 get 请求
func SonyLiv(request *gorequest.SuperAgent) model.Result {
	name := "SonyLiv"
	url := "https://www.sonyliv.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp1, body1, errs1 := request.Get(url).Retry(2, 5).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	if strings.Contains(body1, "geolocation_notsupported") {
		return model.Result{Name: name, Status: model.StatusNo + " (Unavailable)"}
	}
	jwtToken := parseSonyLivToken(body1)

	resp2, body2, errs2 := request.Get("https://apiv2.sonyliv.com/AGL/1.4/A/ENG/WEB/ALL/USER/ULD").
		Set("accept", "application/json, text/plain, */*").
		Set("referer", "https://www.sonyliv.com/").
		Set("device_id", "25a417c3b5f246a393fadb022adc82d5-1715309762699").
		Set("app_version", "3.5.59").
		Set("security_token", jwtToken).End()
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

	resp3, body3, errs3 := request.Get("https://apiv2.sonyliv.com/AGL/3.8/A/ENG/WEB/"+region+
		"/ALL/CONTENT/VIDEOURL/VOD/1000273613/prefetch").
		Set("upgrade-insecure-requests", "1").
		Set("accept", "application/json, text/plain, */*").
		Set("origin", "https://www.sonyliv.com").
		Set("referer", "https://www.sonyliv.com/").
		Set("device_id", "25a417c3b5f246a393fadb022adc82d5-1715309762699").
		Set("security_token", jwtToken).End()
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
		return model.Result{Name: name, Status: model.StatusNo + " (Proxy Detected)", Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get apiv2.sonyliv.com failed with code: %d", resp3.StatusCode)}
}
