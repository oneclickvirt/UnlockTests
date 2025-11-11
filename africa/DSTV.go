package africa

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// DSTV
// authentication.dstv.com 仅 ipv4 且 get 请求
func DSTV(c *http.Client) model.Result {
	name := "DSTV"
	hostname := "dstv.com"
	if c == nil {
		return model.Result{Name: name}
	}
	nowUrl := "https://now.dstv.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(nowUrl)
	if err == nil {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 451:
			return model.Result{Name: name, Status: model.StatusNo}
		case 200:
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		case 403:
			return model.Result{Name: name, Status: model.StatusNo}
		case 404:
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		}
	}
	url := "https://authentication.dstv.com/favicon.ico"
	resp, err = client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 404 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get authentication.dstv.com failed with code: %d", resp.StatusCode)}
}
