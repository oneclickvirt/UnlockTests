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
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	//fmt.Println(body)
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
