package us

import (
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// HBOMax
// www.hbomax.com 仅 ipv4 且 get 请求 (重定向至于 www.max.com 了)
// www.hbonow.com 仅 ipv4 且 get 请求 (重定向至于 www.max.com 了)
func HBOMax(c *http.Client) model.Result {
	name := "HBO Max"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.max.com/"
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
	if strings.Contains(body, "geo-availability") || strings.Contains(resp.Header.Get("location"), "geo-availability") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	t := strings.Split(resp.Header.Get("location"), "/")
	region := ""
	if len(t) >= 4 {
		region = strings.Split(resp.Header.Get("location"), "/")[3]
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToUpper(region)}
}
