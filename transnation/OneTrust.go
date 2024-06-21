package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
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
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	country := utils.ReParse(body, `"country"\s*:\s*"([^"]+)"`)
	stateName := utils.ReParse(body, `"stateName"\s*:\s*"([^"]+)"`)
	if country == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if stateName == "" {
		stateName = "Unknown"
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: country + " " + stateName}
}
