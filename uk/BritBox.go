package uk

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BritBox
// www.britbox.com 双栈 get 请求
func BritBox(request *gorequest.SuperAgent) model.Result {
	name := "BritBox"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.britbox.com/"
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "locationnotsupported") || strings.Contains(body, "locationnotvalidated") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.britbox.com failed with code: %d", resp.StatusCode)}

}
