package transnation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func parseBingRegion(responseBody string) string {
	re := regexp.MustCompile(`Region:"([^"]*)"`)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// Bing
// www.bing.com 双栈 且 post 请求
func Bing(request *gorequest.SuperAgent) model.Result {
	name := "Bing Region"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.bing.com/search?q=www.spiritysdx.top"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		region := parseBingRegion(body)
		if region == "CN" {
			return model.Result{Name: name, Status: model.StatusNo, Region: "cn"}
		}
		if region != "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
		}
	}
	if strings.Contains(body, "cn.bing.com") {
		return model.Result{Name: name, Status: model.StatusNo, Region: "cn"}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.bing.com failed with code: %d", resp.StatusCode)}
}
