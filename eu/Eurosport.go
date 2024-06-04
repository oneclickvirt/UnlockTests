package eu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Eurosport
// www.eurosport.com 双栈
func Eurosport(c *http.Client) model.Result {
	name := "Eurosport RO"
	if c == nil {
		return model.Result{Name: name}
	}
	fakeUuid, _ := uuid.NewV4()
	url := "https://eu3-prod-direct.eurosport.ro/token?realm=eurosport"
	headers := map[string]string{
		"User-Agent":         model.UA_Browser,
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.9",
		"origin":             "https://www.eurosport.ro",
		"referer":            "https://www.eurosport.ro/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"x-device-info":      fmt.Sprintf("escom/0.295.1 (unknown/unknown; Windows/10; %s)", fakeUuid),
		"x-disco-client":     "WEB:UNKNOWN:escom:0.295.1",
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
				Realm string `json:"realm"`
				Token string `json:"token"`
			} `json:"attributes"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body1), &res1); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res1.Data.Attributes.Token != "" {
		//fmt.Println(res1.Data.Attributes.Token)
		sourceSystemId := "eurosport-vid2133403"
		playbackUrl := fmt.Sprintf("https://eu3-prod-direct.eurosport.ro/playback/v2/videoPlaybackInfo/sourceSystemId/%s?usePreAuth=true", sourceSystemId)
		resp2, body2, errs2 := request.Get(playbackUrl).
			Set("Authorization", fmt.Sprintf("Bearer %s", res1.Data.Attributes.Token)).
			End()
		if len(errs2) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
		}
		defer resp2.Body.Close()
		//fmt.Println(body2)
		isBlocked := strings.Contains(body2, "access.denied.geoblocked")
		isOK := strings.Contains(body2, "eurosport-vod")
		if (!isBlocked && !isOK) || isBlocked {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if isOK {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get eu3-prod-direct.eurosport.ro failed with code: %d", resp1.StatusCode)}
}
