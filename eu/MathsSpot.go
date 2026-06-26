package eu

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func MathsSpot(c *http.Client) model.Result {
	name := "MathsSpot"
	hostname := "mathsspot.com"
	if c == nil {
		return model.Result{Name: name}
	}

	headers := map[string]string{
		"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language": "en-US,en;q=0.9",
		"User-Agent":      model.UA_Browser,
	}
	client := utils.SetReqHeaders(utils.Req(c), headers)
	resp, err := client.R().Get("https://mathsspot.com/")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if strings.Contains(body, "FailureServiceNotInRegion") {
		return model.Result{Name: name, Status: model.StatusNo}
	}

	region := utils.ReParse(body, `"countryCode"\s*:\s*"([^"]+)"`)
	if region == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), UnlockType: unlockType}
}
