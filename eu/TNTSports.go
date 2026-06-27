package eu

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// TNTSports
// www.tntsports.co.uk ipv4 get request
func TNTSports(c *http.Client) model.Result {
	name := "TNTSports"
	hostname := "tntsports.co.uk"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://www.tntsports.co.uk/")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if resp.StatusCode == http.StatusTemporaryRedirect &&
		resp.Header.Get("Location") == "https://www.tntsports.co.uk/geoblocking.shtml" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
		}
		re := regexp.MustCompile(`\\"countryCode\\":\\"([A-Z]{2})\\"`)
		if matches := re.FindStringSubmatch(string(b)); len(matches) >= 2 {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(matches[1]), UnlockType: unlockType}
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
