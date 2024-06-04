package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// Popcornflix
// popcornflix-prod.cloud.seachange.com 仅 ipv4 且 get 请求
func Popcornflix(c *http.Client) model.Result {
	name := "Popcornflix"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://popcornflix-prod.cloud.seachange.com/cms/popcornflix/clientconfiguration/versions/2"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get popcornflix-prod.cloud.seachange.com failed with code: %d", resp.StatusCode)}
}
