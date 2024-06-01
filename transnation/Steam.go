package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
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
	for _, c := range resp.Request.Cookies() {
		if c.Name == "steamCountry" {
			i := strings.Index(c.Value, "%")
			if i == -1 {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(c.Value[:i])}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get store.steampowered.com failed with code: %d", resp.StatusCode)}
}
