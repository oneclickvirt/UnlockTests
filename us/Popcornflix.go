package us

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Popcornflix
// popcornflix-prod.cloud.seachange.com 仅 ipv4 且 get 请求
func Popcornflix(c *http.Client) model.Result {
	name := "Popcornflix"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://popcornflix-prod.cloud.seachange.com/cms/popcornflix/clientconfiguration/versions/2"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 400 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get popcornflix-prod.cloud.seachange.com failed with code: %d", resp.StatusCode)}
}
