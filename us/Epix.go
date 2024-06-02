package us

import (
	"encoding/json"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Epix
// api.epix.com 仅 ipv4 且 post 请求
func Epix(request *gorequest.SuperAgent) model.Result {
	name := "MGM+"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api.epix.com/v2/sessions"
	payload := `{"device":{"guid":"7a0baaaf-384c-45cd-a21d-310ca5d3002a","format":"console","os":"web","display_width":1865,"display_height":942,"app_version":"1.0.2","model":"browser","manufacturer":"google"},"apikey":"53e208a9bbaee479903f43b39d7301f7"}`
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Post(url).
		Set("Content-Type", "application/json").
		Set("Connection", "keep-alive").
		Set("traceparent", "00-000000000000000015b7efdb572b7bf2-4aefaea90903bd1f-01").
		Set("sec-ch-ua-mobile", "?0").
		Set("x-datadog-sampling-priority", "1").
		Set("x-datadog-trace-id", "1564983120873880562").
		Set("x-datadog-parent-id", "5399726519264460063").
		Set("Origin", "https://www.mgmplus.com").
		Set("Referer", "https://www.mgmplus.com/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Send(payload).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "error code") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(body, "blocked") {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	var res struct {
		DeviceSession struct {
			SessionToken string `json:"session_token"`
		} `json:"device_session"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	url2 := "https://api.epix.com/v2/movies/16921/play"
	resp2, body2, errs2 := request.Post(url2).
		Set("Content-Type", "application/json").
		Set("X-Session-Token", res.DeviceSession.SessionToken).
		Send("{}").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	var res2 struct {
		Movie struct {
			Entitlements struct {
				Status string `json:"status"`
			} `json:"entitlements"`
		} `json:"movie"`
	}
	if err := json.Unmarshal([]byte(body2), &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	switch res2.Movie.Entitlements.Status {
	case "PROXY_DETECTED":
		return model.Result{Name: name, Status: model.StatusNo, Info: "Proxy Detected"}
	case "GEO_BLOCKED":
		return model.Result{Name: name, Status: model.StatusNo, Info: "Unavailable"}
	case "NOT_SUBSCRIBED":
		return model.Result{Name: name, Status: model.StatusYes}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
}
