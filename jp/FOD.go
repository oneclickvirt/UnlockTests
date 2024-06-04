package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// FOD
// geocontrol1.stream.ne.jp 仅 ipv4 且 get 请求
func FOD(c *http.Client) model.Result {
	name := "FOD(Fuji TV)"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://geocontrol1.stream.ne.jp/fod-geo/check.xml?time=1624504256"
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
	//fmt.Println(body)
	if strings.Contains(body, "FLAG TYPE=\"false\"") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || strings.Contains(body, "true") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get geocontrol1.stream.ne.jp failed with code: %d", resp.StatusCode)}
}
