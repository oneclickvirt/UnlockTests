package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Reddit
// www.reddit.com 仅 ipv4 且 get 请求
func Reddit(c *http.Client) model.Result {
	name := "Reddit"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.reddit.com/"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
