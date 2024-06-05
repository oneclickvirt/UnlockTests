package jp

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// NETRIDE
// trial.net-ride.com 双栈 get 请求
func NETRIDE(c *http.Client) model.Result {
	name := "NETRIDE"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "http://trial.net-ride.com/free/free_dl.php?R_sm_code=456&R_km_url=cabb"
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
	body := string(b)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 || strings.Contains(body, "302 Found") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get trial.net-ride.com failed with code: %d", resp.StatusCode)}
}
