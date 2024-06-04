package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Bing
// www.bing.com 双栈 且 post 请求
func Bing(c *http.Client) model.Result {
	name := "Bing Region"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.bing.com/search?q=www.spiritysdx.top"
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
	if resp.StatusCode == 200 {
		region := utils.ReParse(body, `Region:"([^"]*)"`)
		if region == "CN" {
			return model.Result{Name: name, Status: model.StatusNo, Region: "cn"}
		}
		if region != "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
		}
	}
	if strings.Contains(body, "cn.bing.com") {
		return model.Result{Name: name, Status: model.StatusNo, Region: "cn"}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.bing.com failed with code: %d", resp.StatusCode)}
}
