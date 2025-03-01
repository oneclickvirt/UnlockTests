package eu

import (
	"fmt"
	"net/http"

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
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
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
	if resp.StatusCode == 302 && resp.Header.Get("Location") == "/region-blocked" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// 进一步检查 Viaplay 主站
	if resp.StatusCode == 200 {
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
		if resp2.StatusCode == 302 {
			region := utils.ReParse(resp2.Header.Get("Location"), `/([a-z]{2})/`)
			if region == "" {
				region = utils.ReParse(resp2.Header.Get("Location"), `viaplay.([a-z]{2})`)
			}
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: region}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	// 未知状态码，返回错误
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get checkout.viaplay.pl failed with code: %d", resp.StatusCode)}
}
