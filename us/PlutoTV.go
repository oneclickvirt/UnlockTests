package us

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// PlutoTV
// pluto.tv 仅 ipv4 且 get 请求
func PlutoTV(c *http.Client) model.Result {
	name := "Pluto TV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://pluto.tv/"
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
	if strings.Contains(body, "thanks-for-watching") || strings.Contains(body, "plutotv-is-not-available") ||
		strings.Contains(resp.Request.URL.String(), "plutotv-is-not-available") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusUnexpected, Info: "Rate Limit"}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
