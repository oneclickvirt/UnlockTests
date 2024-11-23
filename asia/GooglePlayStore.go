package asia

import (
	"fmt"
	"io"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// GooglePlayStore
// play.google.com 仅 ipv4 且 get 请求
func GooglePlayStore(c *http.Client) model.Result {
	name := "Google Play Store"
	hostname := "play.google.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://play.google.com/"
	client := utils.Req(c)
	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "en-US;q=0.9",
		"sec-ch-ua":                 `"Chromium";v="131", "Not_A Brand";v="24", "Google Chrome";v="131"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("无法解析响应内容")}
	}
	// 使用正则表达式提取区域信息
	// 匹配<div class="yVZQTb">标签中的内容，直到下一个<字符
	body := string(b)
	matches := utils.ReParse(body, `<div class="yVZQTb">([^<(]+)`) // 应该是地址
	// 检查是否找到匹配内容
	if len(matches) < 2 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{
		Name:       name,
		Status:     model.StatusYes,
		UnlockType: unlockType,
	}
}
