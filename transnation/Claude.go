package transnation

import (
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Claude 检测
func Claude(c *http.Client) model.Result {
	name := "Claude"
	hostname := "claude.ai"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://claude.ai/"
	headers1 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	client1 := utils.Req(c)
	client1 = utils.SetReqHeaders(client1, headers1)
	resp1, err1 := client1.R().Get(url1)
	if err1 == nil && resp1 != nil && resp1.StatusCode == 200 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	url2 := "https://claude.ai/cdn-cgi/trace"
	client2 := utils.Req(c)
	resp2, err2 := client2.R().Get(url2)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
	}
	defer resp2.Body.Close()
	b, err := io.ReadAll(resp2.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	s := string(b)
	lines := strings.Split(s, "\n")
	var location string
	for _, line := range lines {
		if strings.HasPrefix(line, "loc=") {
			location = strings.TrimPrefix(line, "loc=")
			break
		}
	}
	if location == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	if location == "T1" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "TOR"}
	}
	loc := strings.ToLower(location)
	exit := utils.GetRegion(loc, model.ClaudeSupportCountry)
	if exit {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: loc, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo, Region: loc}
}
