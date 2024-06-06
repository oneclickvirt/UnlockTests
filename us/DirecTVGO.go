package us

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// DirecTVGO
// www.directvgo.com 仅 ipv4 且 get 请求
func DirecTVGO(c *http.Client) model.Result {
	name := "DirecTV Go"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.directvgo.com/registrarse"
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
	if strings.Contains(body, "proximamente") || resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		parts := strings.Split(body, "/")
		if len(parts) >= 4 {
			region := parts[3]
			region = strings.ToUpper(region)
			return model.Result{Name: name, Status: model.StatusYes, Region: region}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.directvgo.com failed with code: %d", resp.StatusCode)}
}
