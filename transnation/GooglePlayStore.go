package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// 提取Google Play Store的区域信息
func extractGooglePlayStoreRegion(responseBody string) string {
	// 尝试两种模式匹配区域信息
	patterns := []string{
		`"zQmIje"\s*:\s*"([^"]+)"`,
		`<div class="yVZQTb">([^<(]+)`,
	}
	for _, pattern := range patterns {
		if result := utils.ReParse(responseBody, pattern); result != "" {
			return strings.TrimSpace(result)
		}
	}
	return ""
}

// GooglePlayStore 检测函数
// play.google.com 仅支持 ipv4 且使用 get 请求
func GooglePlayStore(c *http.Client) model.Result {
	name := "Google Play Store"
	hostname := "play.google.com"
	if c == nil {
		return model.Result{Name: name}
	}
	// 设置请求配置
	url := "https://play.google.com/store/games"
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
	}
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("无法解析响应内容")}
	}
	body := string(b)
	if resp.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	region := extractGooglePlayStoreRegion(body)
	if region != "" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		if strings.ToUpper(region) == "CN" {
			return model.Result{
				Name:       name,
				Status:     model.StatusNo,
				Region:     "cn",
				UnlockType: unlockType,
			}
		}
		return model.Result{
			Name:       name,
			Status:     model.StatusYes,
			Region:     strings.ToLower(region),
			UnlockType: unlockType,
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
