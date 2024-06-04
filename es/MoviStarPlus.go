package es

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// MoviStarPlus
// contratar.movistarplus.es 仅 ipv4 且 get 请求
func MoviStarPlus(c *http.Client) model.Result {
	name := "Movistar+"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://contratar.movistarplus.es/"
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
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get contratar.movistarplus.es failed with code: %d", resp.StatusCode)}
}
