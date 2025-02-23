package transnation

import (
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Apple 检测
func Apple(c *http.Client) model.Result {
	name := "Apple"
	hostname := "apple.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://gspe1-ssl.ls.apple.com/pep/gcc"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	location := string(b)
	if location == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	location = utils.TwoToThreeCode(location)
	loc := strings.ToLower(location)
	exit := utils.GetRegion(loc, model.AppleSupportCountry)
	if exit {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: loc, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo, Region: loc}
}
