package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Reddit
// www.reddit.com 仅 ipv4 且 get 请求
func Reddit(request *gorequest.SuperAgent) model.Result {
	name := "Reddit"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.reddit.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 || resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if resp.StatusCode == 403 && strings.Contains(body, "blocked") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.reddit.com failed with code: %d", resp.StatusCode)}
}
