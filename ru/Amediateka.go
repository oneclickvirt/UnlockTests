package ru

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Amediateka
// www.amediateka.ru 仅 ipv4 且 get 请求
func Amediateka(c *http.Client) model.Result {
	name := "Amediateka"
	hostname := "amediateka.ru"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.amediateka.ru/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	// req 自动跟随重定向；非俄用户被跳转到 /unavailable/，通过最终 URL 检测
	if strings.Contains(body, "VPN") || resp.StatusCode == 451 || resp.StatusCode == 455 || resp.StatusCode == 503 ||
		strings.Contains(resp.Request.URL.String(), "/unavailable/") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.amediateka.ru failed with code: %d", resp.StatusCode)}
}
