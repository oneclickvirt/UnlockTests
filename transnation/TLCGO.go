package transnation

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// OneTrust
// geolocation.onetrust.com 双栈 get 请求
func OneTrust(request *gorequest.SuperAgent) model.Result {
	name := "TLC GO"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://geolocation.onetrust.com/cookieconsentpub/v1/geo/location/dnsfeed"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "\"country\":\"US\"") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
