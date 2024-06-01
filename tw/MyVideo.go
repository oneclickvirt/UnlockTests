package tw

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
	"time"
)

// MyVideo
// www.myvideo.net.tw 仅 ipv4 且 get 请求
func MyVideo(request *gorequest.SuperAgent) model.Result {
	name := "MyVideo"
	url := "https://www.myvideo.net.tw/login.do"
	resp, body, errs := request.Timeout(15 * time.Second).Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "serviceAreaBlock") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else {
		return model.Result{Name: name, Status: model.StatusYes}
	}
}