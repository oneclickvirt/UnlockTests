package us

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// AMCPlus
// www.amcplus.com 双栈 且 get 请求
func AMCPlus(c *http.Client) model.Result {
	name := "AMC+"
	hostname := "amcplus.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.amcplus.com/"
	client := utils.Req(c)
	resp1, err1 := client.R().Get(url)
	if err1 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err1}
	}
	defer resp1.Body.Close()
	b, err := io.ReadAll(resp1.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if strings.Contains(body, "Not yet available in your country") || resp1.StatusCode == 403 ||
		resp1.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp1.StatusCode == 302 {
		resp2, err2 := client.R().Get(resp1.Header.Get("Location"))
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == 301 {
			if resp2.Header.Get("Location") == "https://www.amcplus.com/pages/geographic-restriction" {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			resp3, err3 := client.R().Get(resp2.Header.Get("Location"))
			if err3 != nil {
				return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err3}
			}
			defer resp3.Body.Close()
			if resp3.StatusCode == 200 {
				region := utils.ReParse(resp1.Header.Get("Location"), `https://www\.amcplus\.com/countries/(\w{2})`)
				if region != "" {
					result1, result2, result3 := utils.CheckDNS(hostname)
					unlockType := utils.GetUnlockType(result1, result2, result3)
					return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
				}
			}
		}
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get www.amcplus.com failed with code: %d", resp1.StatusCode)}
	}
	if resp1.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: "us", UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get acorn.tv failed with code: %d", resp1.StatusCode)}
}
