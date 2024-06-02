package asia

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// HBOGO
// api2.hbogoasia.com 仅 ipv4 且 get 请求
func HBOGO(request *gorequest.SuperAgent) model.Result {
	name := "HBO GO Asia"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api2.hbogoasia.com/v1/geog?lang=undefined&version=0&bundleId=www.hbogoasia.com"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var hboRes struct {
		Country   string `json:"country"`
		Territory string `json:"territory"`
	}
	//fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &hboRes); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if hboRes.Territory == "" {
		// 解析不到为空则识别为不解锁
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(hboRes.Country)}
}
