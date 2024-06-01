package in

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Zee5
// www.zee5.com 仅 ipv4 且 get 请求
func Zee5(request *gorequest.SuperAgent) model.Result {
	name := "Zee5"
	url := "https://www.zee5.com/"
	request = request.Set("User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0").
		Set("Upgrade-Insecure-Requests", "1")
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if strings.Contains(body, "Access Denied") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		if strings.Contains(resp.Request.URL.String(), "global") {
			return model.Result{Name: name, Status: model.StatusYes, Region: "Global"}
		}
		if resp.Request.URL.String() == url {
			return model.Result{Name: name, Status: model.StatusYes, Region: "in"}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.zee5.com failed with code: %d", resp.StatusCode)}
}
