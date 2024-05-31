package asia

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Bilibili
// B站主体请求逻辑
func Bilibili(request *gorequest.SuperAgent, url string) model.Result {
	name := "Bilibili"
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var res struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "抱歉您所在地区不可观看") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if strings.Contains(body, "抱歉您所在地区不可观看") || strings.Contains(body, "The area is inaccessible") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Code == 0 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.bilibili.com failed with code: %d", resp.StatusCode)}
}
