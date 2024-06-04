package hk

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// MyTvSuper
// www.mytvsuper.com 仅 ipv4 且 get 请求
func MyTvSuper(c *http.Client) model.Result {
	name := "MyTVSuper"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.mytvsuper.com/api/auth/getSession/self/"
	headers := map[string]string{
		"User-Agent":   model.UA_Browser,
		"Content-Type": "application/json",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
