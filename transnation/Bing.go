package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Bing
// www.bing.com 双栈 且 get 请求
func Bing(c *http.Client) model.Result {
	name := "Bing Region"
	if c == nil {
		return model.Result{Name: name}
	}

	url := "https://www.bing.com/search?q=www.spiritysdx.top"
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

	risk_status := strings.Contains(body, "sj_cook.set(\"SRCHHPGUSR\",\"HV\"")

	if resp.StatusCode == 200 {
		region := utils.ReParse(body, `Region:"([^"]*)"`)
		if region == "CN" && strings.Contains(body, "cn.bing.com") {
			info := "Only cn.bing.com"
			if risk_status {
				info += " and Risky"
			}
			return model.Result{Name: name, Status: model.StatusYes, Region: "cn", Info: info}
		}
		if region != "" {
			info := ""
			if risk_status {
				info = "Risky"
			}
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), Info: info}
		}
	}

	if strings.Contains(body, "cn.bing.com") {
		info := "Only cn.bing.com"
		if risk_status {
			info += " and Risky"
		}
		return model.Result{Name: name, Status: model.StatusYes, Region: "cn", Info: info}
	}

	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		info := ""
		if risk_status {
			info = "Risky"
		}
		return model.Result{Name: name, Status: model.StatusBanned, Info: info}
	}

	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.bing.com failed with code: %d", resp.StatusCode)}
}
