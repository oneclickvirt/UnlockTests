package jp

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// PrettyDerby
// api-umamusume.cygames.jp 双栈 且 get 请求
// 有问题 stream error: stream ID 1; INTERNAL_ERROR; received from peer
func PrettyDerby(c *http.Client) model.Result {
	name := "Pretty Derby Japan"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api-umamusume.cygames.jp/"
	headers := map[string]string{
		"User-Agent": model.UA_Dalvik,
	}
	client := utils.ReqDefault(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api-umamusume.cygames.jp failed with code: %d", resp.StatusCode)}
}
