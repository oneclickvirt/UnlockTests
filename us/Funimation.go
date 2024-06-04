package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// Funimation
// www.crunchyroll.com 仅 ipv4 且 get 请求 ( www.funimation.com 重定向为 www.crunchyroll.com 了)
func Funimation(c *http.Client) model.Result {
	name := "Funimation"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.crunchyroll.com/"
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
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	for _, ck := range resp.Request.Cookies() {
		if ck.Name == "region" {
			return model.Result{Name: name, Status: model.StatusYes, Region: ck.Value}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.crunchyroll.com failed with code: %d", resp.StatusCode)}
}
