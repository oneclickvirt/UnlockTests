package ch

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// SkyCh
// sky.ch 双栈 且 get 请求
func SkyCh(c *http.Client) model.Result {
	name := "SKY CH"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://sky.ch/"
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
	if strings.Contains(body, "out-of-country") || strings.Contains(body, "Are you using a VPN") ||
		strings.Contains(body, "Are you using a Proxy or similar Anonymizer technics") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sky.ch failed with code: %d", resp.StatusCode)}
}
