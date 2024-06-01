package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// DirecTVGO
// www.directvgo.com 仅 ipv4 且 get 请求
func DirecTVGO(request *gorequest.SuperAgent) model.Result {
	name := "DirecTV Go"
	url := "https://www.directvgo.com/registrarse"
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "proximamente") || resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		parts := strings.Split(body, "/")
		if len(parts) >= 4 {
			region := parts[3]
			region = strings.ToUpper(region)
			return model.Result{Name: name, Status: model.StatusYes, Region: region}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.directvgo.com failed with code: %d", resp.StatusCode)}
}
