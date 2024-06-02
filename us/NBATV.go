package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NBATV
// www.nba.com 双栈 get 请求
func NBATV(request *gorequest.SuperAgent) model.Result {
	name := "NBA TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.nba.com/watch/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Service is not available in your region") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.nba.com failed with code: %d", resp.StatusCode)}
}
