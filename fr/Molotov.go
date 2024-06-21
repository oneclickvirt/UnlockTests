package fr

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
)

// Molotov
// fapi.molotov.tv 双栈 且 get 请求
func Molotov(c *http.Client) model.Result {
	name := "Molotov"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://fapi.molotov.tv/v1/open-europe/is-france"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	var res struct {
		IsFrance bool `json:"is_france"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
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
