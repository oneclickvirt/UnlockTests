package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// KayoSports
// kayosports.com.au 实际使用 cf 检测，非澳洲请求将一直超时
func KayoSports(c *http.Client) model.Result {
	name := "Kayo Sports"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://kayosports.com.au"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		// if strings.Contains(errs[0].Error(), "i/o timeout") {
		return model.Result{Name: name, Status: model.StatusNo}
		// }
		// return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(resp.Header.Get("Set-Cookie"), "geoblocked=true") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get kayosports.com.au failed with code: %d", resp.StatusCode)}
}
