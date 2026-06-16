package us

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Hulu
// www.hulu.com 仅 ipv4 且 post 请求
func Hulu(c *http.Client) model.Result {
	name := "Hulu"
	hostname := "hulu.com"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"User-Agent":                model.UA_Browser,
		"Accept-Encoding":           "gzip, deflate, br",
		"Cache-Control":             "no-cache",
		"DNT":                       "1",
		"Pragma":                    "no-cache",
		"Sec-CH-UA":                 `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`,
		"Sec-CH-UA-Mobile":          "?0",
		"Sec-CH-UA-Platform":        "Windows",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get("https://www.hulu.com")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if resp.StatusCode == 403 || resp.StatusCode == 451 || strings.Contains(body, "GEO_BLOCKED") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 406 || resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get www.hulu.com failed with code: %d", resp.StatusCode)}
	}

	headers2 := map[string]string{
		"User-Agent":         model.UA_Browser,
		"Accept":             "application/json",
		"Accept-Language":    "en-US,en;q=0.9",
		"Content-Type":       "application/x-www-form-urlencoded; charset=utf-8",
		"Origin":             "https://www.hulu.com",
		"Referer":            "https://www.hulu.com/welcome",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-site",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
	}
	payload := "user_email=&password=&scenario=web_password_login"
	resp, body, err = utils.PostJson(c, "https://auth.hulu.com/v4/web/password/authenticate", payload, headers2)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if strings.Contains(body, "GEO_BLOCKED") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 406 || resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if strings.Contains(body, "LOGIN_FORBIDDEN") || strings.Contains(body, "LOGIN_BAD_REQUEST") ||
		strings.Contains(body, "Your login is invalid. Please refresh the page.") {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
