package us

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// TLCGO
// us1-prod-direct.tlc.com 双栈 get 请求
func TLCGO(c *http.Client) model.Result {
	name := "TLC GO"
	if c == nil {
		return model.Result{Name: name}
	}
	fakeDeviceId, _ := uuid.NewV4()
	url := fmt.Sprintf("https://us1-prod-direct.tlc.com/token?deviceId=%s&realm=go&shortlived=true", fakeDeviceId)
	headers := map[string]string{
		"User-Agent":         model.UA_Browser,
		"accept-language":    "en-US,en;q=0.9",
		"origin":             "https://go.tlc.com",
		"referer":            "https://go.tlc.com/",
		"sec-ch-ua":          "Your_SEC_CH_UA_Here",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"x-device-info":      fmt.Sprintf("tlc/3.17.0 (desktop/desktop; Windows/NT 10.0; %s)", fakeDeviceId),
		"x-disco-client":     "WEB:UNKNOWN:tlc:3.17.0",
		"x-disco-params":     "realm=go,siteLookupKey=tlc,bid=tlc,hn=go.tlc.com,hth=us,features=ar",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
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
		//fmt.Println(res1.Data.Attributes.Token)
		headers2 := map[string]string{
			"User-Agent":         model.UA_Browser,
			"accept-language":    "en-US,en;q=0.9",
			"Authorization":      fmt.Sprintf("Bearer %s", res1.Data.Attributes.Token),
			"origin":             "https://go.tlc.com",
			"referer":            "https://go.tlc.com/",
			"sec-ch-ua":          "Your_SEC_CH_UA_Here",
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": "Windows",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-site",
			"x-disco-client":     "WEB:UNKNOWN:tlc:3.17.0",
			"x-disco-params":     "realm=go,siteLookupKey=tlc,bid=tlc,hn=go.tlc.com,hth=us,features=ar",
		}
		request2 := utils.Gorequest(c)
		request2 = utils.SetGoRequestHeaders(request2, headers2)
		resp2, body2, errs2 := request2.Get("https://us1-prod-direct.tlc.com/cms/routes/home?include=default&decorators=viewingHistory,isFavorite,playbackAllowed&page[items.number]=1&page[items.size]=8").End()
		if len(errs2) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
		}
		defer resp2.Body.Close()
		//fmt.Println(body2)
		isBlocked := strings.Contains(body2, "is not yet available")
		isOK := strings.Contains(body2, "Episodes")
		region := utils.ReParse(body2, `"mainTerritoryCode"\s{0,}:\s{0,}"([^"]+)"`)
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
