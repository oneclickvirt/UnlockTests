package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Gemini
// gemini.google.com 双栈 且 get 请求
func Gemini(c *http.Client) model.Result {
	name := "Gemini"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://gemini.google.com"
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
	if resp.StatusCode == 200 {
		status := false
		if strings.Contains(body, "45631641,null,true") || strings.Contains(body, "45617354,null,true") {
			status = true
		}
		region := utils.ReParse(body, `,2,1,200,"([A-Z]{3})"`)
		region = utils.ThreeToTwoCode(strings.ToLower(region))
		exit := utils.GetRegion(region, model.GeminiSupportCountry)
		result1, result2, result3 := utils.CheckDNS("gemini.google.com")
		unlockType := utils.GetUnlockType(result1, result2, result3)
		if region != "" && exit {
			return model.Result{Name: name, Status: model.StatusYes, Region: region, UnlockType: unlockType}
		} else if status {
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusNetworkErr}
	}
	return model.Result{Name: name, Status: model.StatusNo}
	// return model.Result{Name: name, Status: model.StatusUnexpected,
	// 	Err: fmt.Errorf("get gemini.google.com failed with code: %d", resp.StatusCode)}
}
