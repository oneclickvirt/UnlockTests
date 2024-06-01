package nl

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// VideoLand
// api.videoland.com 双栈 且 post 请求
func VideoLand(request *gorequest.SuperAgent) model.Result {
	name := "Videoland"
	payload := `{"operationName":"IsOnboardingGeoBlocked","variables":{},"query":"query IsOnboardingGeoBlocked {\n  isOnboardingGeoBlocked\n}\n"}`
	resp, body, errs := utils.PostJson(request, "https://api.videoland.com/subscribe/videoland-account/graphql", payload,
		map[string]string{"connection": "keep-alive"},
		map[string]string{"apollographql-client-name": "apollo_accounts_base"},
		map[string]string{"traceparent": "00-cab2dbd109bf1e003903ec43eb4c067d-623ef8e56174b85a-01"},
		map[string]string{"origin": "https://www.videoland.com"},
		map[string]string{"referer": "https://www.videoland.com/"},
		map[string]string{"accept": "application/json, text/plain, */*"},
	)
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Data struct {
			Blocked bool `json:"isOnboardingGeoBlocked"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Data.Blocked {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
