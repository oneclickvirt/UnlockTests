package us

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Starz
// www.starz.com 双栈 get 请求
func Starz(request *gorequest.SuperAgent) model.Result {
	name := "Starz"
	if request == nil {
		return model.Result{Name: name}
	}
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("Referer", "https://www.starz.com/us/en/")
	// client.Headers.Set("Authtokenauthorization", "")
	url := "https://www.starz.com/sapi/header/v1/starz/us/09b397fc9eb64d5080687fc8a218775b" // 请求有tls校验
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	authorization := string(b)
	// fmt.Printf(authorization)
	if authorization != "" && !strings.Contains(authorization, "AccessDenied") {
		resp2, body2, errs2 := request.Get("https://auth.starz.com/api/v4/User/geolocation").
			Set("AuthTokenAuthorization", authorization).
			Retry(2, 5).End()
		if len(errs2) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
		}
		defer resp2.Body.Close()
		var res struct {
			IsAllowedAccess  bool   `json:"isAllowedAccess"`
			IsAllowedCountry bool   `json:"isAllowedCountry"`
			IsKnownProxy     bool   `json:"isKnownProxy"`
			Country          string `json:"country"`
		}
		// fmt.Println(body2)
		if err := json.Unmarshal([]byte(body2), &res); err != nil {
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res.IsAllowedAccess && res.IsAllowedCountry && !res.IsKnownProxy {
			return model.Result{Name: name, Status: model.StatusYes}
		}
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.starz.com failed with code: %d", resp.StatusCode)}
}
