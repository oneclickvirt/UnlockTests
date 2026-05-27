package jp

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// NETRIDE
// trial.net-ride.com 双栈 get 请求
func NETRIDE(c *http.Client) model.Result {
	name := "NETRIDE"
	hostname := "net-ride.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "http://trial.net-ride.com/free/free_dl.php?R_sm_code=456&R_km_url=cabb"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	body := string(b)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// req 自动跟随重定向；日本 IP 被跳转到内容页（最终 URL 与原始 URL 不同）
	if strings.Contains(body, "302 Found") ||
		(resp.StatusCode == 200 && resp.Response != nil && resp.Response.Request != nil &&
			resp.Response.Request.URL.String() != url) {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get trial.net-ride.com failed with code: %d", resp.StatusCode)}
}
