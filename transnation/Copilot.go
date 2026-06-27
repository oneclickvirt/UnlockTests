package transnation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Copilot
// copilot.microsoft.com dual-stack get request
func Copilot(c *http.Client) model.Result {
	name := "Microsoft Copilot"
	hostname := "copilot.microsoft.com"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://copilot.microsoft.com/c/api/user?api-version=2")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusFound:
		if resp.Header.Get("Location") == "/" {
			return model.Result{Name: name, Status: model.StatusBanned}
		}
		return model.Result{Name: name, Status: model.StatusNo}
	case http.StatusForbidden:
		return model.Result{Name: name, Status: model.StatusNo}
	case http.StatusOK:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
		}
		var res struct {
			RegionCode string `json:"regionCode"`
		}
		if err := json.Unmarshal(b, &res); err != nil {
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res.RegionCode != "" {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(res.RegionCode), UnlockType: unlockType}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get copilot.microsoft.com failed with code: %d", resp.StatusCode)}
}
