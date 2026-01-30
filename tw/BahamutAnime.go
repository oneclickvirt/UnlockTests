package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// BahamutAnime
// ani.gamer.com.tw 仅 ipv4 且 get 请求
func BahamutAnime(c *http.Client) model.Result {
	name := "Bahamut Anime"
	hostname := "gamer.com.tw"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	// 获取 device ID
	resp1, err := client.R().
		SetHeader("x-custom-headers", "true").
		Get("https://ani.gamer.com.tw/ajax/getdeviceid.php")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp1.Body.Close()
	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	var res1 struct {
		AnimeSn  int    `json:"animeSn"`
		Deviceid string `json:"deviceid"`
	}
	if err := json.Unmarshal(b1, &res1); err != nil {
		body := string(b1)
		// 检测 Cloudflare 封禁或系统异常
		if strings.Contains(body, "Just a moment") || strings.Contains(body, "系統異常回報") {
			return model.Result{Name: name, Status: model.StatusNo, Info: "Banned by cloudflare"}
		}
		// 如果是 HTML 响应（非 JSON），判定为不可用
		if strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	// 测试第一个视频 (sn=37783)
	resp2, err := client.R().
		SetHeader("x-custom-headers", "true").
		Get("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=37783&device=" + res1.Deviceid)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body2 := string(b2)
	var res2 struct {
		AnimeSn int `json:"animeSn"`
	}
	if err := json.Unmarshal(b2, &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	// 如果第一个视频可以访问
	if res2.AnimeSn != 0 {
		// 测试第二个视频 (sn=38832) - 仅限台湾地区
		resp3, err := client.R().
			SetHeader("x-custom-headers", "true").
			Get("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=38832&device=" + res1.Deviceid)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		defer resp3.Body.Close()
		b3, err := io.ReadAll(resp3.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		var res3 struct {
			AnimeSn int `json:"animeSn"`
		}
		if err := json.Unmarshal(b3, &res3); err != nil {
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		// 如果第二个视频也可以访问，说明是台湾地区
		if res3.AnimeSn != 0 {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, Region: "tw", UnlockType: unlockType}
		}
		// 第二个视频不可访问，获取实际地区
		resp4, err := client.R().Get("https://ani.gamer.com.tw/cdn-cgi/trace")
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		defer resp4.Body.Close()
		b4, err := io.ReadAll(resp4.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		bodyString := string(b4)
		index := strings.Index(bodyString, "loc=")
		if index == -1 {
			return model.Result{Name: name, Status: model.StatusUnexpected}
		}
		bodyString = bodyString[index+4:]
		index = strings.Index(bodyString, "\n")
		if index == -1 {
			return model.Result{Name: name, Status: model.StatusUnexpected}
		}
		loc := bodyString[:index]
		if len(loc) == 2 {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(loc), UnlockType: unlockType}
		}
	}
	// 检查是否有地区限制的错误消息
	if strings.Contains(body2, "\u5f88\u62b1\u6b49\uff01\u672c\u7bc0\u76ee\u56e0\u6388\u6b0a\u56e0\u7d20\u7121\u6cd5\u5728\u60a8\u7684\u6240\u5728\u5340\u57df\u64ad\u653e\u3002") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// 获取主页检查 data-geo 标记
	resp5, err := client.R().Get("https://ani.gamer.com.tw/")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp5.Body.Close()
	b5, err := io.ReadAll(resp5.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body5 := string(b5)
	// 如果包含 animeSn 或设备验证异常，且主页有 data-geo 标记
	if (strings.Contains(body2, "animeSn") ||
		strings.Contains(body2, "\u88dd\u7f6e\u9a57\u8b49\u7570\u5e38\uff01")) && strings.Contains(body5, "data-geo") {
		var location string
		resp6, err := client.R().Get("https://ani.gamer.com.tw/cdn-cgi/trace")
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		defer resp6.Body.Close()
		b6, err := io.ReadAll(resp6.Body)
		if err != nil {
			return utils.HandleNetworkError(c, hostname, err, name)
		}
		bodyString := string(b6)
		index := strings.Index(bodyString, "loc=")
		if index != -1 {
			bodyString = bodyString[index+4:]
			index = strings.Index(bodyString, "\n")
			if index != -1 {
				loc := bodyString[:index]
				if len(loc) == 2 {
					location = strings.ToLower(loc)
				}
			}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: location, UnlockType: unlockType}
	} else if resp2.StatusCode == 403 || resp2.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ani.gamer.com.tw failed with code: %d", resp1.StatusCode)}
}