package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// DAnimeStore
// animestore.docomo.ne.jp 仅 ipv4 且 get 请求
func DAnimeStore(c *http.Client) model.Result {
	name := "D Anime Store"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://animestore.docomo.ne.jp/animestore/reg_pc"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 || strings.Contains(body, "海外") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 && body != "" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get animestore.docomo.ne.jp failed with code: %d", resp.StatusCode)}
}
