package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ParamountPlus
// www.paramountplus.com 双栈 且 get 请求
func ParamountPlus(request *gorequest.SuperAgent) model.Result {
	name := "Paramount+"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.paramountplus.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 && strings.Contains(body, "\"country_name_intl\":\"International\"") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.paramountplus.com failed with code: %d", resp.StatusCode)}
}
