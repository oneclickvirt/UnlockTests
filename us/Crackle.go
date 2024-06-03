package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Crackle
// prod-api.crackle.com 双栈 get 请求 有问题
// {"path":"/appconfig","version":"v2.0.0","status":"400","timestamp":"2024-05-31T10:28:34.542Z","error":{"message":"Platform Key is not specified","type":"ApiError","code":121,"details":{}}}
func Crackle(request *gorequest.SuperAgent) model.Result {
	name := "Crackle"
	url := "https://prod-api.crackle.com/appconfig"
	request = request.Set("User-Agent", model.UA_Browser).
		Set("Accept-Language", "en-US,en;q=0.9").
		Set("Content-Type", "application/json").
		Set("Origin", "https://www.crackle.com").
		Set("Referer", "https://www.crackle.com/").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "same-site").
		Set("sec-ch-ua", "${UA_SEC_CH_UA}").
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"Windows\"").
		Set("x-crackle-apiversion", "v2.0.0").
		Set("x-crackle-brand", "crackle").
		Set("x-crackle-platform", "5FE67CCA-069A-42C6-A20F-4B47A8054D46")
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	fmt.Println(body)
	// TODO 获取地区
	// x-crackle-region
	if strings.Contains(body, "302 Found") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get prod-api.crackle.com failed with code: %d", resp.StatusCode)}
}
