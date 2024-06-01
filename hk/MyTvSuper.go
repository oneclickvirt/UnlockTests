package hk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// MyTvSuper
// www.mytvsuper.com 仅 ipv4 且 get 请求
func MyTvSuper(request *gorequest.SuperAgent) model.Result {
	name := "MyTVSuper"
	url := "https://www.mytvsuper.com/api/auth/getSession/self/"
	request = request.Set("User-Agent", model.UA_Browser).Set("Content-Type", "application/json")
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var mytvsuperRes struct {
		Region      int    `json:"region"`
		CountryCode string `json:"country_code"`
	}
	if err := json.Unmarshal([]byte(body), &mytvsuperRes); err != nil {
		if strings.Contains(body, "HK") {
			return model.Result{Name: name, Status: model.StatusYes}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if mytvsuperRes.Region == 1 && mytvsuperRes.CountryCode == "HK" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if mytvsuperRes.Region != 1 || mytvsuperRes.CountryCode != "HK" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.mytvsuper.com failed with code: %d", resp.StatusCode)}
}
