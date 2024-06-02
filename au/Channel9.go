package au

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Channel9
// login.nine.com.au 双栈 且 get 请求
func Channel9(request *gorequest.SuperAgent) model.Result {
	name := "Channel 9"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://login.nine.com.au"
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get login.nine.com.au failed with code: %d", resp.StatusCode)}
}
