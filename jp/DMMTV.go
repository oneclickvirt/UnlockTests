package jp

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// DMMTV
// api.beacon.dmm.com 双栈 且 post 请求
func DMMTV(request *gorequest.SuperAgent) model.Result {
	name := "DMM TV"
	resp, bodyBytes, errs := utils.PostJson(request, "https://api.beacon.dmm.com/v1/streaming/start",
		`{"player_name":"dmmtv_browser","player_version":"0.0.0","content_type_detail":"VOD_SVOD","content_id":"11uvjcm4fw2wdu7drtd1epnvz","purchase_product_id":null}`,
		map[string]string{"User-Agent": model.UA_Browser},
	)
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		IsBkocked   bool   `json:"is_blocked"`
		BlockStatus string `json:"block_status"`
	}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		if strings.Contains(string(bodyBytes), "UNAUTHORIZED") {
			return model.Result{Name: name, Status: model.StatusYes}
		}
		if strings.Contains(string(bodyBytes), "FOREIGN") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.IsBkocked {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if !res.IsBkocked {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.beacon.dmm.com failed with code: %d", resp.StatusCode)}
}
