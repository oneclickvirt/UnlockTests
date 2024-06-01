package transnation

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Steam
// store.steampowered.com 仅 ipv4 且 get 请求
func Steam(request *gorequest.SuperAgent) model.Result {
	name := "Steam Currency"
	url := "https://store.steampowered.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	cookies := resp.Header.Get("Set-Cookie")
	if strings.Contains(cookies, "steamCountry=") {
		region := strings.ToLower(strings.ReplaceAll(cookies, "steamCountry=", "")[0:2])
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
