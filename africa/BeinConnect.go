package africa

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

// BeinConnect
// proxies.bein-mena-production.eu-west-2.tuc.red 仅 ipv4 且 get 请求
func BeinConnect(c *http.Client) model.Result {
	name := "Bein Sports Connect"
	hostname := "beinconnect.com.tr"
	if c == nil {
		return model.Result{Name: name}
	}
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	url := "https://proxies.bein-mena-production.eu-west-2.tuc.red/proxy/availableOffers"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("BeinConnect Get request failed: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("BeinConnect read resp.Body failed: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	//fmt.Println(body)
	if strings.Contains(body, "Unavailable For Legal Reasons") ||
		resp.StatusCode == 403 || resp.StatusCode == 451 {
		if model.EnableLoger {
			Logger.Info("BeinConnect access denied due to legal reasons: " + body)
		}
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 500 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	if model.EnableLoger {
		Logger.Info(fmt.Sprintf("BeinConnect unexpected response code: %d", resp.StatusCode))
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get proxies.bein-mena-production.eu-west-2.tuc.red failed with code: %d", resp.StatusCode)}
}
