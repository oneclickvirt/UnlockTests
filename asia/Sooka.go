package asia

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Sooka
// sooka.my ipv4 get request
func Sooka(c *http.Client) model.Result {
	name := "Sooka"
	hostname := "sooka.my"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://sooka.my/")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	case http.StatusForbidden:
		return model.Result{Name: name, Status: model.StatusNo}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get sooka.my failed with code: %d", resp.StatusCode)}
	}
}
