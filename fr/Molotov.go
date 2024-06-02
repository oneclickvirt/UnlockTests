package fr

import (
	"encoding/json"
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Molotov
// fapi.molotov.tv 双栈 且 get 请求
func Molotov(request *gorequest.SuperAgent) model.Result {
	name := "Molotov"
	if request == nil {
		return model.Result{Name: name}
	}
	resp, body, errs := request.Get("https://fapi.molotov.tv/v1/open-europe/is-france").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		IsFrance bool `json:"is_france"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.IsFrance {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if !res.IsFrance {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get fapi.molotov.tv failed with code: %d", resp.StatusCode)}
}
