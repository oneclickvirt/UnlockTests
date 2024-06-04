package transnation

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Spotify
// spclient.wg.spotify.com 双栈 且 post 请求
func Spotify(c *http.Client) model.Result {
	name := "Spotify Registration"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://spclient.wg.spotify.com/signup/public/v1/account"
	resp, body, errs := utils.PostJson(c, url, "birth_day=11&birth_month=11&birth_year=2000&collect_personal_info=undefined&creation_flow=&creation_point=https%3A%2F%2Fwww.spotify.com%2Fhk-en%2F&displayname=Gay%20Lord&gender=male&iagree=1&key=a1e486e2729f46d6bb368d6b2bcda326&platform=www&referrer=&send-email=0&thirdpartyemail=0&identifier_token=AgE6YTvEzkReHNfJpO114514",
		map[string]string{"Accept-Language": "en",
			"User-Agent":    model.UA_Browser,
			"content-type":  "application/json",
			"cache-control": "no-cache"})
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Status            int    `json:"status"`
		Country           string `json:"country"`
		IsCountryLaunched bool   `json:"is_country_launched"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Status == 320 || res.Status == 120 || resp.StatusCode == 401 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Status == 311 && res.IsCountryLaunched {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(res.Country)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get spclient.wg.spotify.com failed with code: %d", resp.StatusCode)}
}
