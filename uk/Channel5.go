package uk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

const defaultChannel5Auth = "0_rZDiY0hp_TNcDyk2uD-Kl40HqDbXs7hOawxyqPnbI"

// Channel5
// cassie.channel5.com 仅 ipv4 且 get 请求
func Channel5(c *http.Client) model.Result {
	name := "Channel 5"
	hostname := "channel5.com"
	if c == nil {
		return model.Result{Name: name}
	}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	auth := strings.TrimSpace(os.Getenv("UNLOCKTESTS_CHANNEL5_AUTH"))
	if auth == "" {
		auth = defaultChannel5Auth
	}
	url := fmt.Sprintf("https://cassie.channel5.com/api/v2/live_media/my5desktopng/C5.json?timestamp=%d&auth=%s", timestamp, url.QueryEscape(auth))
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	// fmt.Println(body)
	var res struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		if err.Error() == `invalid character '<' looking for beginning of value` {
			return model.Result{Name: name, Status: model.StatusBanned}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Code == "3000" || strings.Contains(body, "this service is only available in restricted regions") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Code == "4003" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get cassie.channel5.com failed with code: %d", resp.StatusCode)}
}
