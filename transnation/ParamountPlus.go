package transnation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func ParamountPlus(c *http.Client) model.Result {
	name := "ParamountPlus"
	hostname := "www.paramountplus.com"
	if c == nil {
		return model.Result{Name: name}
	}

	client := utils.Req(c)
	resp, err := client.R().Get("https://www.paramountplus.com/")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnavailableForLegalReasons || resp.StatusCode == http.StatusNotFound {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(string(b), "intl") {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
