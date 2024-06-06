package asia

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// HotStar
// api.hotstar.com 双栈 get 请求
func HotStar(c *http.Client) model.Result {
	name := "HotStar"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.hotstar.com/o/v1/page/1557?offset=0&size=20&tao=0&tas=20"
	request := utils.Gorequest(c)
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request = utils.SetGoRequestHeaders(request, headers)
	resp, _, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 475 || resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	resp1, _, errs1 := request.Get("https://www.hotstar.com").Retry(2, 5).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get api.hotstar.com failed with code: %d %d", resp.StatusCode, resp1.StatusCode)}
	}
	defer resp1.Body.Close()
	if resp1.StatusCode == 301 || resp.StatusCode == 475 || resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	//fmt.Println(body)
	u := resp1.Header.Get("Location")
	if u == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	// fmt.Println(u)
	t := strings.SplitN(u, "/", 4)
	if len(t) < 4 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.ToLower(t[3]) == "us" {
		return model.Result{Name: name, Status: model.StatusNo, Region: t[3]}
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: t[3]}
	// return model.Result{Name: name, Status: model.StatusUnexpected,
	// 	Err: fmt.Errorf("get api.hotstar.com failed with code: %d", resp.StatusCode)}
}
