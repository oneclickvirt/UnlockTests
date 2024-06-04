package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Shudder
// www.shudder.com 双栈 get 请求
func Shudder(c *http.Client) model.Result {
	name := "Shudder"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.shudder.com/"
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
	if strings.Contains(body, "not available") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "movies") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.shudder.com failed with code: %d", resp.StatusCode)}
}
