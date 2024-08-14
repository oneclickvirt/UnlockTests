package hk

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// HoyTV
// hoytv-live-stream.hoy.tv 双栈 且 post 请求
func HoyTV(c *http.Client) model.Result {
	name := "Hoy TV"
	hostname := "hoy.tv"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://hoytv-live-stream.hoy.tv/ch78/index-fhd.m3u8"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	} else if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get hoytv-live-stream.hoy.tv failed with code: %d", resp.StatusCode)}
}
