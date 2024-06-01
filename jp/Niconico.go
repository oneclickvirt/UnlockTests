package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Niconico
// www.nicovideo.jp 仅 ipv4 且 get 请求
func Niconico(request *gorequest.SuperAgent) model.Result {
	name := "Niconico"
	url1 := "https://www.nicovideo.jp/watch/so40278367" // 进击的巨人
	//url2 := "https://www.nicovideo.jp/watch/so23017073" // 假面骑士
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url1).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "同じ地域") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.nicovideo.jp failed with code: %d", resp.StatusCode)}
}
