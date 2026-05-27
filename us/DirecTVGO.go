package us

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// DirecTVGO
// www.directvgo.com 仅 ipv4 且 get 请求
func DirecTVGO(c *http.Client) model.Result {
	name := "DirecTV Go"
	hostname := "directvgo.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.directvgo.com/registrarse"
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
	if strings.Contains(body, "proximamente") || resp.StatusCode == 403 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		// req 自动跟随重定向，从最终 URL 提取地区（如 directvgo.com/mx/）
		region := ""
		if resp.Response != nil && resp.Response.Request != nil {
			finalURL := resp.Response.Request.URL.String()
			region = utils.ReParse(finalURL, `directvgo\.com/([a-z]{2})/`)
			region = strings.ToUpper(region)
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		if region != "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
		}
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.directvgo.com failed with code: %d", resp.StatusCode)}
}
