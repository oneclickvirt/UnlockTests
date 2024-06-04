package au

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Channel10
// 10play.com.au 仅 ipv4 且 get 请求
// https://e410fasadvz.global.ssl.fastly.net/geo 仅 ipv4 且 get 请求
// https://10play.com.au/geo-web 仅 ipv4 且 get 请求
func Channel10(c *http.Client) model.Result {
	name := "Channel 10"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://10play.com.au/geo-web"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Sorry, 10 play is not available in your region.") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	url = "https://e410fasadvz.global.ssl.fastly.net/geo"
	resp, body, errs = request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	var res struct {
		State string `json:"state"`
		Allow bool   `json:"allow"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "not available") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if !res.Allow {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if res.Allow && res.State != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: res.State}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get 10play.com.au failed with code: %d", resp.StatusCode)}
}
