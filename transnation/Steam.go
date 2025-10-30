package transnation

import (
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Steam
// store.steampowered.com 仅 ipv4 且 get 请求
func Steam(c *http.Client) model.Result {
	name := "Steam Store"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://store.steampowered.com/"
	url2 := "https://steamcommunity.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	cookies := resp.Header.Get("Set-Cookie")
	if strings.Contains(cookies, "steamCountry=") {
		region := strings.ToLower(strings.ReplaceAll(cookies, "steamCountry=", "")[0:2])
		resp2, err2 := client.R().Get(url2)
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region),
				Info: "Community Unavailable"}
		} else {
			defer resp2.Body.Close()
		}
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), Info: "Community Available"}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
