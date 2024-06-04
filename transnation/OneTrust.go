package transnation

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// OneTrust
// geolocation.onetrust.com 双栈 get 请求
func OneTrust(c *http.Client) model.Result {
	name := "OneTrust"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://geolocation.onetrust.com/cookieconsentpub/v1/geo/location/dnsfeed"
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
	country := utils.ReParse(body, `"country"\s*:\s*"([^"]+)"`)
	stateName := utils.ReParse(body, `"stateName"\s*:\s*"([^"]+)"`)
	if strings.ToLower(country) == "us" {
		return model.Result{Name: name, Status: model.StatusYes, Region: country + " " + stateName}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
