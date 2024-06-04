package us

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// FuboTV
// api.fubo.tv 仅 ipv4 且 get 请求
func FuboTV(c *http.Client) model.Result {
	name := "Fubo TV"
	if c == nil {
		return model.Result{Name: name}
	}
	randNum := strconv.Itoa(rand.Intn(2))
	url := "https://api.fubo.tv/appconfig/v1/homepage?platform=web&client_version=R20230310." + randNum + "&nav=v0"
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
	if strings.Contains(body, "No Subscription") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if strings.Contains(body, "Forbidden IP") {
		return model.Result{Name: name, Status: model.StatusYes, Info: "IP Forbidden"}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
