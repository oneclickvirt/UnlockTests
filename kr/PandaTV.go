package kr

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// PandaTV
// api.pandalive.co.kr 仅 ipv4 且 get 请求
func PandaTV(c *http.Client) model.Result {
	name := "PandaTV"
	hostname := "pandalive.co.kr"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.pandalive.co.kr/v1/live/play"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Result{Status: model.StatusNetworkErr, Err: err}
	// }
	switch resp.StatusCode {
	case 403:
		return model.Result{Name: name, Status: model.StatusNo}
	case 400:
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected code: %d", resp.StatusCode)}
	}
}
