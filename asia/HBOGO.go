package asia

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// HBOGO
// api2.hbogoasia.com 仅 ipv4 且 get 请求
func HBOGO(c *http.Client) model.Result {
	name := "HBO GO Asia"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api2.hbogoasia.com/v1/geog?lang=undefined&version=0&bundleId=www.hbogoasia.com"
	request := utils.Gorequest(c)
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
