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
	if strings.Contains(body, "Not yet available in your country") || resp.StatusCode == 403 ||
		resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == 302 {
		resp2, err2 := client.R().Get(resp.Header.Get("Location"))
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		defer resp2.Body.Close()
		if resp2.StatusCode == 301 && resp2.Header.Get("Location") == "https://www.amcplus.com/pages/geographic-restriction" {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get www.amcplus.com failed with code: %d", resp.StatusCode)}
	}
	if resp.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get acorn.tv failed with code: %d", resp.StatusCode)}
}
