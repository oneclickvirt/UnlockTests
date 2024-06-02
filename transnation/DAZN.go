package transnation

import (
	"encoding/json"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// DAZN
// startup.core.indazn.com 仅 ipv4 且 post 请求
func DAZN(request *gorequest.SuperAgent) model.Result {
	name := "Dazn"
	if request == nil {
		return model.Result{Name: name}
	}
	resp, bodyBytes, errs := utils.PostJson(request, "https://startup.core.indazn.com/misl/v5/Startup",
		`{"LandingPageKey":"generic","Languages":"zh-CN,zh,en","Platform":"web","PlatformAttributes":{},"Manufacturer":"","PromoCode":"","Version":"2"}`,
		map[string]string{"User-Agent": model.UA_Browser},
	)
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
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
	if err := json.Unmarshal(bodyBytes, &daznRes); err != nil {
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
