package tw

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// MyVideo
// 结合重定向检查和页面内容分析，检测 www.myvideo.net.tw 的访问权限
func MyVideo(c *http.Client) model.Result {
	const (
		name     = "MyVideo"
		hostname = "myvideo.net.tw"
		url      = "https://www.myvideo.net.tw/login.do"
	)
	if c == nil {
		return model.Result{Name: name}
	}
	// 发起 GET 请求（req 会自动跟随重定向，通过最终 URL 判断结果）
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{
			Name:   name,
			Status: model.StatusNetworkErr,
			Err:    err,
		}
	}
	defer resp.Body.Close()
	// 检查最终跳转 URL
	if resp.Response != nil && resp.Response.Request != nil {
		finalPath := resp.Response.Request.URL.Path
		if strings.Contains(finalPath, "serviceAreaBlock") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if strings.Contains(finalPath, "goLoginPage") {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		}
	}
	// 如果最终 URL 无法判断，读取响应体并检查内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{
			Name:   name,
			Status: model.StatusNetworkErr,
			Err:    fmt.Errorf("unable to parse response body"),
		}
	}
	body := string(bodyBytes)
	// 根据页面内容判断区域限制
	if strings.Contains(body, "serviceAreaBlock") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// 如果页面未显示限制，检查 DNS 并返回解锁状态
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{
		Name:       name,
		Status:     model.StatusYes,
		UnlockType: unlockType,
	}
}
