package ch

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// SkyCh
// sky.ch 双栈 且 get 请求
func SkyCh(c *http.Client) model.Result {
	name := "SKY CH"
	if c == nil {
		return model.Result{Name: name}
	}
	hostname := "sky.ch"
	client := utils.Req(c)
	url := "https://gateway.prd.sky.ch/user/customer/create"
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		body := string(b)
		if body == `{"message": "", "code": "GEO_BLOCKED"}` {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == 405 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	originalUrl := "https://sky.ch/"
	resp2, err := client.R().Get(originalUrl)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp2.Body.Close()
	b, err := io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)

	if strings.Contains(body, "out-of-country") || strings.Contains(body, "Are you using a VPN") ||
		strings.Contains(body, "Are you using a Proxy or similar Anonymizer technics") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp2.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sky.ch failed with code: %d", resp2.StatusCode)}
}
