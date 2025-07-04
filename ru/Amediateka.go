package ru

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
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
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if strings.Contains(body, "VPN") || resp.StatusCode == 451 || resp.StatusCode == 455 || resp.StatusCode == 503 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 301 && resp.Header.Get("Location") == "https://www.amediateka.ru/unavailable/index.html" {
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
