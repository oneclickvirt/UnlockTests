package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"time"
)

// WorldFlipper
// api.worldflipper.jp 双栈 且 get 请求
func WorldFlipper(c *http.Client) model.Result {
	name := "World Flipper Japan"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.worldflipper.jp/"
	headers := map[string]string{
		"User-Agent": model.UA_Dalvik,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).Timeout(10*time.Second).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 404 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.worldflipper.jp failed with code: %d", resp.StatusCode)}
}
