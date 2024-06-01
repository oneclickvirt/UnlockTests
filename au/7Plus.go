package au

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Au7plus
// 7plus.com.au 仅 ipv4 且 get 请求
// 7plus-sevennetwork.akamaized.net 有问题 - 无论如何请求都失败
func Au7plus(request *gorequest.SuperAgent) model.Result {
	name := "7plus"
	url := "https://7plus-sevennetwork.akamaized.net/media/v1/dash/live/cenc/5303576322001/68dca38b-85d7-4dae-b1c5-c88acc58d51c/f4ea4711-514e-4cad-824f-e0c87db0a614/225ec0a0-ef18-4b7c-8fd6-8dcdd16cf03a/1x/segment0.m4f?akamai_token=exp=1672500385~acl=/media/v1/dash/live/cenc/5303576322001/68dca38b-85d7-4dae-b1c5-c88acc58d51c/f4ea4711-514e-4cad-824f-e0c87db0a614/*~hmac=800e1e1d1943addf12b71339277c637c7211582fe12d148e486ae40d6549dbde"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// fmt.Println(body)
	// fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		resp1, _, errs1 := request.Get("https://7plus.com.au/").End()
		if len(errs1) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
		}
		defer resp1.Body.Close()
		// fmt.Println(body)
		if resp1.StatusCode == 403 || resp1.StatusCode == 451 {
			return model.Result{Name: name, Status: model.StatusNo}
		} else if resp1.StatusCode == 200 {
			return model.Result{Name: name, Status: model.StatusYes}
		} else {
			return model.Result{Name: name, Status: model.StatusUnexpected,
				Err: fmt.Errorf("get 7plus.com.au failed with code: %d %d", resp.StatusCode, resp1.StatusCode)}
		}
	}
}
