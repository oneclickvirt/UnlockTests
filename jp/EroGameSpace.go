package jp

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// EroGameSpace
// erogamescape.org 仅支持 IPv4 且仅接受 GET 请求
func EroGameSpace(c *http.Client) model.Result {
	const (
		name     = "EroGameSpace"
		hostname = "erogamescape.org"
		url      = "https://erogamescape.org/"
	)
	// 发起请求
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
	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{
			Name:   name,
			Status: model.StatusNetworkErr,
			Err:    err,
		}
	}
	body := string(bodyBytes)
	// 根据 HTTP 响应状态码判断结果
	switch resp.StatusCode {
	case 403, 451:
		return model.Result{Name: name, Status: model.StatusNo}
	case 200:
		// 检查响应内容是否包含 "18歳" 字样
		if strings.Contains(body, "18歳") {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{
				Name:       name,
				Status:     model.StatusYes,
				UnlockType: unlockType,
			}
		}
	}
	// 返回意外结果的错误
	return model.Result{
		Name:   name,
		Status: model.StatusUnexpected,
		Err:    fmt.Errorf("unexpected response from %s with status code: %d", hostname, resp.StatusCode),
	}
}
