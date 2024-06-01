package ch

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// SkyCh
// sky.ch 双栈 且 get 请求
func SkyCh(request *gorequest.SuperAgent) model.Result {
	name := "SKY CH"
	url := "https://sky.ch/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if strings.Contains(body, "out-of-country") || strings.Contains(body, "Are you using a VPN") ||
		strings.Contains(body, "Are you using a Proxy or similar Anonymizer technics") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sky.ch failed with code: %d", resp.StatusCode)}
}
