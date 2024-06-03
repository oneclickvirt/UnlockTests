package transnation

import (
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func extractFields(jsonStr string) (string, string) {
	// 定义正则表达式
	countryRegex := regexp.MustCompile(`"country"\s*:\s*"([^"]+)"`)
	stateNameRegex := regexp.MustCompile(`"stateName"\s*:\s*"([^"]+)"`)
	// 查找匹配的部分
	countryMatch := countryRegex.FindStringSubmatch(jsonStr)
	stateNameMatch := stateNameRegex.FindStringSubmatch(jsonStr)
	var country, stateName string
	if len(countryMatch) > 1 {
		country = countryMatch[1]
	}
	if len(stateNameMatch) > 1 {
		stateName = stateNameMatch[1]
	}
	return country, stateName
}

// OneTrust
// geolocation.onetrust.com 双栈 get 请求
func OneTrust(request *gorequest.SuperAgent) model.Result {
	name := "OneTrust"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://geolocation.onetrust.com/cookieconsentpub/v1/geo/location/dnsfeed"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	country, stateName := extractFields(body)
	if strings.ToLower(country) == "us" {
		return model.Result{Name: name, Status: model.StatusYes, Region: country + " " + stateName}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
