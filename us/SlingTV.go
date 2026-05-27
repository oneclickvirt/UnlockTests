package us

import (
	"fmt"
	"net/http"

	req "github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// SlingTV
// www.sling.com 双栈 且 get 请求
func SlingTV(c *http.Client) model.Result {
	name := "Sling TV"
	hostname := "sling.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.sling.com/"
	headers := map[string]string{
		"User-Agent": model.UA_Dalvik,
	}
	client := utils.ReqDefault(c)
	client = utils.SetReqHeaders(client, headers)
	// 禁止自动跟随重定向：sling.com 对非美国 IP 返回 302 作为不可用的信号
	client.SetRedirectPolicy(req.NoRedirectPolicy())
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
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.sling.com failed with code: %d", resp.StatusCode)}
}
