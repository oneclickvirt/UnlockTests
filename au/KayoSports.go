package au

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// KayoSports
// kayosports.com.au 实际使用 cf 检测，非澳洲请求将一直超时
func KayoSports(c *http.Client) model.Result {
	name := "Kayo Sports"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://kayosports.com.au"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	defer resp.Body.Close()
	// fmt.Println(resp.Header.Get("Set-Cookie"))
	if strings.Contains(resp.Header.Get("Set-Cookie"), "geoblocked=true") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get kayosports.com.au failed with code: %d", resp.StatusCode)}
}
