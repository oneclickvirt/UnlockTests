package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// Binge
// auth.streamotion.com.au 双栈 get 请求
func Binge(c *http.Client) model.Result {
	name := "Binge"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://auth.streamotion.com.au"
	request := utils.Gorequest(c)
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get auth.streamotion.com.au failed with code: %d", resp.StatusCode)}

}
