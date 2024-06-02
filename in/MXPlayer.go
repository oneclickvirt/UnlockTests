package in

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// MXPlayer
// www.mxplayer.in 仅 ipv4 且 get 请求
func MXPlayer(request *gorequest.SuperAgent) model.Result {
	name := "MX Player"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.mxplayer.in/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	//fmt.Println(resp.Header.Get("set-cookie"))
	if strings.Contains(body, "We are currently not available in your region") ||
		strings.Contains(body, "403 ERROR") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.Header.Get("set-cookie") != "" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.mxplayer.in failed with code: %d", resp.StatusCode)}
}
