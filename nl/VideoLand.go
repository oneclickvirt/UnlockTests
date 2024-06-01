package nl

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// VideoLand
// api.videoland.com 双栈 且 post 请求
func VideoLand(request *gorequest.SuperAgent) model.Result {
	name := "Videoland"
	url := "https://api.videoland.com/subscribe/videoland-account/graphql"
	payload := `{"operationName":"IsOnboardingGeoBlocked","variables":{},"query":"query IsOnboardingGeoBlocked {\n  isOnboardingGeoBlocked\n}\n"}`
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("connection", "keep-alive")
	client.Headers.Set("apollographql-client-name", "apollo_accounts_base")
	client.Headers.Set("traceparent", "00-cab2dbd109bf1e003903ec43eb4c067d-623ef8e56174b85a-01")
	client.Headers.Set("origin", "https://www.videoland.com")
	client.Headers.Set("referer", "https://www.videoland.com/")
	client.Headers.Set("accept", "application/json, text/plain, */*")
	resp, err := client.R().SetBodyString(payload).Post(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	//fmt.Println(body)
	var res struct {
		Data struct {
			Blocked bool `json:"isOnboardingGeoBlocked"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "\"isOnboardingGeoBlocked\":true") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Data.Blocked {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
