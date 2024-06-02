package jp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Telasa
// api-videopass-anon.kddi-video.com 双栈 get 请求
func Telasa(request *gorequest.SuperAgent) model.Result {
	name := "Telasa"
	if request == nil {
		return model.Result{Name: name}
	}
	resp, body, errs := request.Get("https://api-videopass-anon.kddi-video.com/v1/playback/system_status").
		Set("X-Device-ID", "d36f8e6b-e344-4f5e-9a55-90aeb3403799").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Status struct {
			Type    string `json:"type"`
			Subtype string `json:"subtype"`
		} `json:"status"`
	}
	// fmt.Println(body)
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "RequestForbidden") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Status.Subtype == "IPLocationNotAllowed" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Status.Type != "" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api-videopass-anon.kddi-video.com failed with code: %d", resp.StatusCode)}
}
