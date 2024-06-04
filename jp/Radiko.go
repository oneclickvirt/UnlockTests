package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"regexp"
	"strings"
)

// Radiko
// radiko.jp 仅 ipv4 且 get 请求
func Radiko(c *http.Client) model.Result {
	name := "Radiko"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://radiko.jp/area?_=1625406539531"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "class=\"OUT\"") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(body, "JAPAN") {
		tempList := strings.Split(body, "\n")
		var location string
		for _, line := range tempList {
			if strings.Contains(line, "JAPAN") {
				// 使用 strings.Fields 来分割字符串，并获取第二个字段
				fields := strings.Fields(line)
				if len(fields) < 2 {
					break
				}
				secondField := fields[1]
				// 使用正则表达式删除最后一个 '>' 字符之前的所有内容
				re := regexp.MustCompile(`.*>`)
				location = re.ReplaceAllString(secondField, "")
				break
			}
		}
		if location != "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: location}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get radiko.jp failed with code: %d", resp.StatusCode)}
}
