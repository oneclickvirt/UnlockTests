package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// ATTNOW - DirectvStream
// www.atttvnow.com 双栈 且 get 请求
func DirectvStream(c *http.Client) model.Result {
	name := "Directv Stream"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.atttvnow.com/"
	request := utils.Gorequest(c)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes,
		Err: fmt.Errorf("get www.atttvnow.com failed with code: %d", resp.StatusCode)}
}
