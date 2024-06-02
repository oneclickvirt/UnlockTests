package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// KOCOWA
// www.kocowa.com 仅 ipv4 且 get 请求
func KOCOWA(request *gorequest.SuperAgent) model.Result {
	name := "KOCOWA"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.kocowa.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 403 || strings.Contains(body, "is not available in your region or country") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.kocowa.com failed with code: %d", resp.StatusCode)}
}
