package jp

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// MGStage
// www.mgstage.com 仅 ipv4 且 get 请求
func MGStage(c *http.Client) model.Result {
	name := "MGStage"
	hostname := "mgstage.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.mgstage.com/"
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
	case 200:
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected code: %d", resp.StatusCode)}
	}
}
