package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// OptusSports
// sport.optus.com.au 双栈 get 请求
func OptusSports(request *gorequest.SuperAgent) model.Result {
	name := "Optus Sports"
	url := "https://sport.optus.com.au/api/userauth/validate/web/username/restriction.check@gmail.com"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sport.optus.com.au failed with code: %d", resp.StatusCode)}
}
