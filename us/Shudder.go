package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// Shudder
// www.shudder.com 双栈 get 请求
func Shudder(c *http.Client) model.Result {
	name := "Shudder"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.shudder.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	if strings.Contains(body, "not available") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "movies") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.shudder.com failed with code: %d", resp.StatusCode)}
}
