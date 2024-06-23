package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

// FuboTV
// api.fubo.tv 仅 ipv4 且 get 请求
func FuboTV(c *http.Client) model.Result {
	name := "Fubo TV"
	hostname := "fubo.tv"
	if c == nil {
		return model.Result{Name: name}
	}
	randNum := strconv.Itoa(rand.Intn(2))
	url := "https://api.fubo.tv/appconfig/v1/homepage?platform=web&client_version=R20230310." + randNum + "&nav=v0"
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
	if strings.Contains(body, "No Subscription") {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	} else if strings.Contains(body, "Forbidden IP") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "IP Forbidden"}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
