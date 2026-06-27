package vn

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// GalaxyPlay
// api.glxplay.io ipv4 get request
func GalaxyPlay(c *http.Client) model.Result {
	name := "Galaxy Play"
	hostname := "glxplay.io"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://api.glxplay.io/account/device/new")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	if resp.StatusCode == http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
		}
		body := string(b)
		if strings.HasPrefix(body, "<") ||
			strings.Contains(body, `"errorCode": 495`) ||
			strings.Contains(body, "not available in your region") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.glxplay.io failed with code: %d", resp.StatusCode)}
}
