package asia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

// MolaTV
// mola.tv 仅 ipv4 且 get 请求
func MolaTV(c *http.Client) model.Result {
	name := "Mola TV"
	hostname := "mola.tv"
	if c == nil {
		return model.Result{Name: name}
	}
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	url := "https://mola.tv/api/v2/videos/geoguard/check/vd30491025"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("mola.tv Get request failed: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("mola.tv can not parse body: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	//fmt.Println(body)
	var res struct {
		Data struct {
			Type       string `json:"type"`
			Id         string `json:"id"`
			Attributes struct {
				IsAllowed bool `json:"isAllowed"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		if strings.Contains(body, "\"isAllowed\":false") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if model.EnableLoger {
			Logger.Info("mola.tv can not parse json: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	if !res.Data.Attributes.IsAllowed {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if res.Data.Attributes.IsAllowed {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	if model.EnableLoger {
		Logger.Info(fmt.Sprintf("mola.tv unexpected response code: %d", resp.StatusCode))
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mola.tv failed with code: %d", resp.StatusCode)}
}
