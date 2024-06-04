package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// ABCiView
// api.iview.abc.net.au 仅 ipv4 且 get 请求
func ABCiView(c *http.Client) model.Result {
	name := "ABC iView"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.iview.abc.net.au/v2/show/abc-kids-live-stream/video/LS1604H001S00?embed=highlightVideo,selectedSeries"
	request := utils.Gorequest(c)
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "unavailable outside Australia") || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.iview.abc.net.au failed with code: %d", resp.StatusCode)}
}
