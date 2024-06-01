package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"regexp"
	"strings"
)

func TikTokCountry(body string) string {
	re := regexp.MustCompile(`"region":"(\w+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// TikTok
// www.tiktok.com 仅 ipv4 且 get 请求
func TikTok(request *gorequest.SuperAgent) model.Result {
	name := "TikTok"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get("https://www.tiktok.com/explore").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "https://www.tiktok.com/hk/notfound") {
		return model.Result{Name: name, Status: model.StatusNo, Region: "hk"}
	}
	if region := TikTokCountry(body); region != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.tiktok.com failed with code: %d", resp.StatusCode)}
}
