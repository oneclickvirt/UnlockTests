package eu

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Eurosport
// www.eurosport.com 双栈
func Eurosport(request *gorequest.SuperAgent) model.Result {
	name := "Eurosport RO"
	if request == nil {
		return model.Result{Name: name}
	}
	fakeUuid, _ := uuid.NewV4()
	resp1, body1, errs1 := request.Get("https://eu3-prod-direct.eurosport.ro/token?realm=eurosport").
		Set("accept", "*/*").
		Set("accept-language", "en-US,en;q=0.9").
		Set("origin", "https://www.eurosport.ro").
		Set("referer", "https://www.eurosport.ro/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"Windows\"").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("x-device-info", fmt.Sprintf("escom/0.295.1 (unknown/unknown; Windows/10; %s)", fakeUuid)).
		Set("x-disco-client", "WEB:UNKNOWN:escom:0.295.1").
		Set("User-Agent", model.UA_Browser).
		End()
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
