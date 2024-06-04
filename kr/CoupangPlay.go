package kr

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// CoupangPlay
// www.coupangplay.com 仅 ipv4 且 get 请求
func CoupangPlay(c *http.Client) model.Result {
	name := "Coupang Play"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.coupangplay.com/"
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
	//fmt.Println(resp.Request.URL.String())
	if strings.Contains(body, "is not available in your region") ||
		strings.Contains(resp.Request.URL.String(), "not-available") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.coupangplay.com failed with code: %d", resp.StatusCode)}
}
