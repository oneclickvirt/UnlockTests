package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// KKTV
// api.kktv.me 仅 ipv4 且 get 请求
func KKTV(c *http.Client) model.Result {
	name := "KKTV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.kktv.me/v3/ipcheck"
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
	var res struct {
		Data struct {
			Country   string `json:"country"`
			IsAllowed bool   `json:"is_allowed"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "\"is_allowed\":false") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Data.Country == "TW" && res.Data.IsAllowed {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if !res.Data.IsAllowed {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.kktv.me failed with head: %d", resp.StatusCode)}
}
