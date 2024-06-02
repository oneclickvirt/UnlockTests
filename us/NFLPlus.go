package us

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NFLPlus
// www.dazn.com 仅 ipv4 且 get 请求
// https://www.nfl.com/plus/ 重定向至于 https://nfl.com/dazn-watch-gp-row 约等于仅使用 dazn 进行观看
func NFLPlus(request *gorequest.SuperAgent) model.Result {
	name := "NFL+"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://nfl.com/dazn-watch-gp-row"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "nflgamepass") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get https://nfl.com/dazn-watch-gp-row failed with code: %d", resp.StatusCode)}
}
