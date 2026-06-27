package vn

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// ClipTV
// cliptv.vn ipv4 get request
func ClipTV(c *http.Client) model.Result {
	name := "Clip TV"
	hostname := "cliptv.vn"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://cliptv.vn/truyen-hinh")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
		}
		if strings.Contains(string(b), "Sorry, this video is not available in your country.") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get cliptv.vn failed with code: %d", resp.StatusCode)}
}
