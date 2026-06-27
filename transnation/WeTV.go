package transnation

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// WeTV
// wetv.vip dual-stack get request
func WeTV(c *http.Client) model.Result {
	name := "WeTV"
	hostname := "wetv.vip"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://wetv.vip/")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
		}
		body := string(b)
		if strings.Contains(body, "static.wetvinfo.com/static/vg-cn/tar") {
			return model.Result{Name: name, Status: model.StatusNo, Region: "cn"}
		}
		re := regexp.MustCompile(`(?i)"areaPhoneId"\s*:\s*"\+([0-9]+)"`)
		if match := re.FindStringSubmatch(body); len(match) > 1 {
			if region := utils.CountryCodeToAlpha2(match[1]); region != "" {
				result1, result2, result3 := utils.CheckDNS(hostname)
				unlockType := utils.GetUnlockType(result1, result2, result3)
				return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), UnlockType: unlockType}
			}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get wetv.vip failed with code: %d", resp.StatusCode)}
}
