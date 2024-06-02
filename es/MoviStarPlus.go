package es

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// MoviStarPlus
// contratar.movistarplus.es 仅 ipv4 且 get 请求
func MoviStarPlus(request *gorequest.SuperAgent) model.Result {
	name := "Movistar+"
	if request == nil {
		return model.Result{Name: name}
	}
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get("https://contratar.movistarplus.es/").End()
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
