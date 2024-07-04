package africa

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

// DSTV
// authentication.dstv.com 仅 ipv4 且 get 请求
func DSTV(c *http.Client) model.Result {
	name := "DSTV"
	hostname := "dstv.com"
	if c == nil {
		return model.Result{Name: name}
	}
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	url := "https://authentication.dstv.com/favicon.ico"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("DSTV Get request failed: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 404 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	if model.EnableLoger {
		Logger.Info(fmt.Sprintf("DSTV unexpected response code: %d", resp.StatusCode))
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get authentication.dstv.com failed with code: %d", resp.StatusCode)}
}
