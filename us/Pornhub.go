package us

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Pornhub
// www.pornhub.com 仅 ipv4 且 get 请求
// 美国某些州已经因法律原因限制访问
func Pornhub(c *http.Client) model.Result {
	name := "Pornhub"
	hostname := "pornhub.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.pornhub.com/"
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
	if strings.Contains(body, "has blocked access") ||
		strings.Contains(body, "not available in your region") ||
		strings.Contains(body, "age verification") && strings.Contains(body, "required by law") ||
		strings.Contains(body, "restricted in your location") ||
		strings.Contains(body, "blocked in your state") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Blocked by regional law"}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Access Forbidden"}
	}
	if resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusUnexpected, Info: "Rate Limit"}
	}
	if resp.StatusCode == 200 && (strings.Contains(body, "pornhub") ||
		strings.Contains(body, "Free Porn") ||
		strings.Contains(body, "Premium")) {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("unexpected response with code: %d", resp.StatusCode)}
}
