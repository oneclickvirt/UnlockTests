package tw

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func LineTV(c *http.Client) model.Result {
	name := "LineTV.TW"
	hostname := "linetv.tw"
	if c == nil {
		return model.Result{Name: name}
	}

	client := utils.Req(c)
	resp, err := client.R().Get("https://www.linetv.tw/drama/11829/eps/1")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if !strings.Contains(body, "window.__INITIAL_STATE__") {
		return model.Result{Name: name, Status: model.StatusNo}
	}

	reEps := regexp.MustCompile(`"eps_info"\s*:\s*\[`)
	reDuration := regexp.MustCompile(`"durationInMs"\s*:\s*\d+`)
	if reEps.MatchString(body) && reDuration.MatchString(body) {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
