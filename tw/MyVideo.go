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
	// 设置 HTTP 客户端的重定向行为
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // 禁止自动跟随重定向
	}
	// 发起 GET 请求
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
	// 检查重定向逻辑
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		switch location {
		case "/serviceAreaBlock.do":
			return model.Result{Name: name, Status: model.StatusNo}
		case "/goLoginPage.do":
			return model.Result{Name: name, Status: model.StatusYes}
		default:
			return model.Result{
				Name:   name,
				Status: model.StatusUnexpected,
				Err:    fmt.Errorf("unexpected redirection to: %s", location),
			}
		}
	}
	// 如果未发生重定向，读取响应体并检查内容
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
