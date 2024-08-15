package eu

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Viaplay
// checkout.viaplay.pl 仅 ipv4 且 get 请求
func Viaplay(c *http.Client) model.Result {
	name := "Viaplay"
	hostname := "viaplay.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://checkout.viaplay.pl/?recommended=viaplay"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == 302 && resp.Header.Get("Location") == "/region-blocked" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	// }
	// body := string(b)
	if resp.StatusCode == 200 {
		url2 := "https://viaplay.com/"
		resp2, err2 := client.R().Get(url2)
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == 302 {
			region := utils.ReParse(resp2.Header.Get("Location"), `/([a-z]{2})/`)
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: region}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get checkout.viaplay.pl failed with code: %d", resp.StatusCode)}
}
