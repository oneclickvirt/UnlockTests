package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// ProjectSekai
// game-version.sekai.colorfulpalette.org 仅 ipv4 且 get 请求
func ProjectSekai(c *http.Client) model.Result {
	name := "Project Sekai - Colorful Stage"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://game-version.sekai.colorfulpalette.org/1.8.1/3ed70b6a-8352-4532-b819-108837926ff5"
	headers := map[string]string{
		"User-Agent": model.UA_Pjsekai,
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
		Err: fmt.Errorf("get game-version.sekai.colorfulpalette.org failed with code: %d", resp.StatusCode)}
}
