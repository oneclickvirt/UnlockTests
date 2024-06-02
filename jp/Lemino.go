package jp

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Lemino
// if.lemino.docomo.ne.jp 双栈 且 get 请求
func Lemino(request *gorequest.SuperAgent) model.Result {
	name := "Lemino"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://if.lemino.docomo.ne.jp/v1/user/delivery/watch/ready"
	//request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).
		Set("Accept", "application/json, text/plain, */*").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		Set("Content-Type", "application/json").
		Set("Origin", "https://lemino.docomo.ne.jp").
		Set("Referer", "https://lemino.docomo.ne.jp/").
		Set("Sec-CH-UA-Mobile", "?0").
		Set("Sec-CH-UA-Platform", "\"Windows\"").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "same-site").
		Set("X-Service-Token", "f365771afd91452fa279863f240c233d").
		Set("X-Trace-ID", "556db33f-d739-4a82-84df-dd509a8aa179").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Send("{\"inflow_flows\":[null,\"crid://plala.iptvf.jp/group/b100ce3\"],\"play_type\":1,\"key_download_only\":null,\"quality\":null,\"groupcast\":null,\"avail_status\":\"1\",\"terminal_type\":3,\"test_account\":0,\"content_list\":[{\"kind\":\"main\",\"service_id\":null,\"cid\":\"00lm78dz30\",\"lid\":\"a0lsa6kum1\",\"crid\":\"crid://plala.iptvf.jp/vod/0000000000_00lm78dymn\",\"preview\":0,\"trailer\":0,\"auto_play\":0,\"stop_position\":0}]}").
		Retry(2, 5).End()
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
		Err: fmt.Errorf("get if.lemino.docomo.ne.jp failed with code: %d", resp.StatusCode)}
}
