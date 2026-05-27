package us

import (
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// PeacockTV
// www.peacocktv.com 双栈 get 请求
func PeacockTV(c *http.Client) model.Result {
	name := "Peacock TV"
	hostname := "peacocktv.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.peacocktv.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	// req 自动跟随重定向；非美国用户被跳转到含 "unavailable" 的页面，通过最终 URL 检测
	if strings.Contains(resp.Request.URL.String(), "unavailable") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
}
