package tw

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// MyVideo
// www.myvideo.net.tw 仅 ipv4 且 get 请求
func MyVideo(c *http.Client) model.Result {
	name := "MyVideo"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.myvideo.net.tw/login.do"
	request := utils.Gorequest(c)
	resp, body, errs := request.Get(url).End()
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
