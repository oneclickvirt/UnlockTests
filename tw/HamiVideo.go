package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"time"
)

// HamiVideo
// hamivideo.hinet.net 仅 ipv4 且 get 请求
func HamiVideo(request *gorequest.SuperAgent) model.Result {
	name := "Hami Video"
	url := "https://hamivideo.hinet.net/api/play.do?id=OTT_VOD_0000249064&freeProduct=1"
	request = request.Set("User-Agent", model.UA_Browser).Timeout(15 * time.Second)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Code == "06001-107" {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if res.Code == "06001-106" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get hamivideo.hinet.net failed with code: %d", resp.StatusCode)}
}