package kr

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func SOOP(c *http.Client) model.Result {
	name := "SOOP"
	hostname := "vod.sooplive.co.kr"
	if c == nil {
		return model.Result{Name: name}
	}

	client := utils.Req(c)
	resp, err := client.R().Get("https://vod.sooplive.co.kr/player/97464151")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get vod.sooplive.co.kr failed with code: %d", resp.StatusCode)}
}
