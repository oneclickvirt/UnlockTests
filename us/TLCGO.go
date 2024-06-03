package us

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func extractValue(text, regexPattern string) string {
	re := regexp.MustCompile(regexPattern)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// TLCGO
// us1-prod-direct.tlc.com 双栈 get 请求
func TLCGO(request *gorequest.SuperAgent) model.Result {
	name := "TLC GO"
	if request == nil {
		return model.Result{Name: name}
	}
	fakeDeviceId, _ := uuid.NewV4()
	url := fmt.Sprintf("https://us1-prod-direct.tlc.com/token?deviceId=%s&realm=go&shortlived=true", fakeDeviceId)
	request = request.Set("User-Agent", model.UA_Browser).Retry(2, 5).
		Set("accept-language", "en-US,en;q=0.9").
		Set("origin", "https://go.tlc.com").
		Set("referer", "https://go.tlc.com/").
		Set("sec-ch-ua", "Your_SEC_CH_UA_Here").
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("x-device-info", fmt.Sprintf("tlc/3.17.0 (desktop/desktop; Windows/NT 10.0; %s)", fakeDeviceId)).
		Set("x-disco-client", "WEB:UNKNOWN:tlc:3.17.0").
		Set("x-disco-params", "realm=go,siteLookupKey=tlc,bid=tlc,hn=go.tlc.com,hth=us,features=ar")
	resp1, body1, errs1 := request.Get(url).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	//fmt.Println(body1)
	var res1 struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body1), &res1); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res1.Data.Attributes.Token != "" {
		// fmt.Println(res1.Data.Attributes.Token)
		resp2, body2, errs2 := request.Get("https://us1-prod-direct.tlc.com/cms/routes/home?include=default&decorators=viewingHistory,isFavorite,playbackAllowed&page[items.number]=1&page[items.size]=8").
			Set("accept-language", "en-US,en;q=0.9").
			Set("Authorization", fmt.Sprintf("Bearer %s", res1.Data.Attributes.Token)).
			Set("origin", "https://go.tlc.com").
			Set("referer", "https://go.tlc.com/").
			Set("sec-ch-ua", "Your_SEC_CH_UA_Here").
			Set("sec-ch-ua-mobile", "?0").
			Set("sec-ch-ua-platform", "Windows").
			Set("sec-fetch-dest", "empty").
			Set("sec-fetch-mode", "cors").
			Set("sec-fetch-site", "same-site").
			Set("x-disco-client", "WEB:UNKNOWN:tlc:3.17.0").
			Set("x-disco-params", "realm=go,siteLookupKey=tlc,bid=tlc,hn=go.tlc.com,hth=us,features=ar").
			End()
		if len(errs2) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
		}
		defer resp2.Body.Close()
		//fmt.Println(body2)
		isBlocked := strings.Contains(body2, "is not yet available")
		isOK := strings.Contains(body2, "Episodes")
		region := extractValue(body2, `"mainTerritoryCode"\s{0,}:\s{0,}"([^"]+)"`)
		if !isBlocked && !isOK {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if isBlocked {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if isOK {
			return model.Result{Name: name, Status: model.StatusYes, Region: region}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get us1-prod-direct.tlc.com failed with code: %d", resp1.StatusCode)}
}
