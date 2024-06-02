package us

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// FuboTV
// api.fubo.tv 仅 ipv4 且 get 请求
func FuboTV(request *gorequest.SuperAgent) model.Result {
	name := "Fubo TV"
	if request == nil {
		return model.Result{Name: name}
	}
	randNum := strconv.Itoa(rand.Intn(2))
	url := "https://api.fubo.tv/appconfig/v1/homepage?platform=web&client_version=R20230310." + randNum + "&nav=v0"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "No Subscription") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if strings.Contains(body, "Forbidden IP") {
		return model.Result{Name: name, Status: model.StatusYes, Info: "IP Forbidden"}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
