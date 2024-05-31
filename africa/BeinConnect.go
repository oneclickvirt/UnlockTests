package africa

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// BeinConnect
// proxies.bein-mena-production.eu-west-2.tuc.red 仅 ipv4 且 get 请求
func BeinConnect(request *gorequest.SuperAgent) model.Result {
	name := "Bein Sports Connect"
	url := "https://proxies.bein-mena-production.eu-west-2.tuc.red/proxy/availableOffers"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Unavailable For Legal Reasons") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 500 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get proxies.bein-mena-production.eu-west-2.tuc.red failed with code: %d", resp.StatusCode)}
}
