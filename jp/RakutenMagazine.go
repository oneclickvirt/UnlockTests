package jp

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// RakutenMagazine
// magazine.rakuten.co.jp 仅 ipv4 且 get 请求
func RakutenMagazine(c *http.Client) model.Result {
	name := "Rakuten Magazine"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://data-cloudauthoring.magazine.rakuten.co.jp/rem_repository/////////.key"
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
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get magazine.rakuten.co.jp failed with code: %d", resp.StatusCode)}
}
