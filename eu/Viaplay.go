package eu

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Viaplay
// checkout.viaplay.pl 仅 ipv4 且 get 请求
func Viaplay(c *http.Client) model.Result {
	name := "Viaplay"
	hostname := "viaplay.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://checkout.viaplay.pl/?recommended=viaplay"
	client := utils.Req(c)
	// 发送请求并检查错误
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	// 处理 HTTP 状态码
	if resp.StatusCode == 403 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	// 进一步检查 Viaplay 主站
	if resp.StatusCode == 200 {
		// req 自动跟随重定向；若最终落在 region-blocked 路径则不可用
		if resp.Response != nil && resp.Response.Request != nil &&
			strings.Contains(resp.Response.Request.URL.Path, "region-blocked") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		url2 := "https://viaplay.com/"
		resp2, err2 := client.R().Get(url2)
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
		}
		// 避免空指针
		defer func() {
			if resp2 != nil && resp2.Body != nil {
				resp2.Body.Close()
			}
		}()
		if resp2.StatusCode == 404 {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		// req 自动跟随重定向，从最终 URL 提取地区（如 viaplay.se/viaplay.no 等）
		finalURL2 := ""
		if resp2.Response != nil && resp2.Response.Request != nil {
			finalURL2 = resp2.Response.Request.URL.String()
		}
		region := utils.ReParse(finalURL2, `/([a-z]{2})/`)
		if region == "" {
			region = utils.ReParse(finalURL2, `viaplay\.([a-z]{2})`)
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: region}
	}
	// 未知状态码，返回错误
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get checkout.viaplay.pl failed with code: %d", resp.StatusCode)}
}
