package transnation

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Sora
// sora.com 仅 ipv4 且 get 请求
func Sora(c *http.Client) model.Result {
	name := "Sora"
	if c == nil {
		return model.Result{Name: name}
	}
	// 创建 cookie jar
	jar, _ := cookiejar.New(nil)
	c.Jar = jar
	// 第一次请求获取地区信息
	client := utils.Req(c)
	resp, err := client.R().Get("https://sora.com/cdn-cgi/trace")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	// 解析位置信息
	s := string(b)
	i := strings.Index(s, "loc=")
	if i == -1 {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i == -1 {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	loc := s[:i]
	// 第二次请求检查认证状态
	resp, err = client.R().Get("https://sora.com/backend/authenticate")
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	// 检查各种状态
	if strings.Contains(string(b), "Attention Required") {
		return model.Result{Name: name, Status: model.StatusBanned, Info: "VPN Blocked"}
	}
	lowLoc := strings.ToLower(loc)
	if resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusRestricted, Region: lowLoc, Info: "429 Rate limit"}
	}
	if loc == "T1" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "tor"}
	}
	if exit := utils.GetRegion(lowLoc, model.GptSupportCountry); exit {
		return model.Result{Name: name, Status: model.StatusYes, Region: lowLoc}
	}
	return model.Result{Name: name, Status: model.StatusNo, Region: lowLoc}
}
