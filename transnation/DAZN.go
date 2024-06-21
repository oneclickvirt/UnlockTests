package transnation

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// DAZN
// startup.core.indazn.com 仅 ipv4 且 post 请求
func DAZN(c *http.Client) model.Result {
	name := "Dazn"
	if c == nil {
		return model.Result{Name: name}
	}
	resp, body, err := utils.PostJson(c, "https://startup.core.indazn.com/misl/v5/Startup",
		`{"LandingPageKey":"generic","Languages":"zh-CN,zh,en","Platform":"web","PlatformAttributes":{},"Manufacturer":"","PromoCode":"","Version":"2"}`,
		map[string]string{"User-Agent": model.UA_Browser},
	)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	var daznRes struct {
		Region struct {
			IsAllowed             bool   `json:"isAllowed"`
			DisallowedReason      string `json:"disallowedReason"`
			GeolocatedCountry     string `json:"GeolocatedCountry"`
			GeolocatedCountryName string `json:"GeolocatedCountryName"`
		} `json:"Region"`
	}
	if err := json.Unmarshal([]byte(body), &daznRes); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if daznRes.Region.IsAllowed {
		return model.Result{
			Name: name, Status: model.StatusYes,
			Region: daznRes.Region.GeolocatedCountry,
		}
	}
	return model.Result{
		Name: name, Status: model.StatusNo, Info: daznRes.Region.DisallowedReason,
	}
}
