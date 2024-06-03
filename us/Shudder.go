package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Shudder
// www.shudder.com 双栈 get 请求
func Shudder(request *gorequest.SuperAgent) model.Result {
	name := "Shudder"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.shudder.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "not available") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "movies") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.shudder.com failed with code: %d", resp.StatusCode)}
}
