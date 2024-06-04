package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Steam
// store.steampowered.com 仅 ipv4 且 get 请求
func Steam(c *http.Client) model.Result {
	name := "Steam Currency"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://store.steampowered.com/"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
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
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get store.steampowered.com failed with code: %d", resp.StatusCode)}
}
