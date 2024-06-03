package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// MetaAI
// www.meta.ai 双栈 且 get 请求 有问题
func MetaAI(request *gorequest.SuperAgent) model.Result {
	name := "MetaAI"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.meta.ai/"
	request = request.Set("User-Agent", model.UA_Browser).
		Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Set("Accept-Language", "en-US,en;q=0.9").
		Set("sec-ch-ua", "${UA_SEC_CH_UA}").
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "document").
		Set("sec-fetch-mode", "navigate").
		Set("sec-fetch-site", "none").
		Set("sec-fetch-user", "?1").
		Set("upgrade-insecure-requests", "1")
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// fmt.Println(body)
	// 检查是否被阻止
	if strings.Contains(body, "AbraGeoBlockedErrorRoot") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "AbraGeoBlocked"}
	}
	// 检查是否成功
	if strings.Contains(body, "AbraHomeRootConversationQuery") {
		start := strings.Index(body, `"code"`)
		if start != -1 {
			start = strings.Index(body[start:], `"`) + start + 1
			end := strings.Index(body[start:], `"`) + start
			code := body[start:end]
			region := strings.Split(code, "_")[1]
			if region != "" {
				return model.Result{Name: name, Status: model.StatusYes, Region: region}
			}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.meta.ai failed with code: %d", resp.StatusCode)}
}
