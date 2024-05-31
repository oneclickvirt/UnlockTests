package asia

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Hotstar
// api.hotstar.com 双栈 get 请求
func Hotstar(request *gorequest.SuperAgent) model.Result {
	name := "Hotstar"
	url := "https://api.hotstar.com/o/v1/page/1557?offset=0&size=20&tao=0&tas=20"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 475 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 401 {
		resp, _, errs = request.Get("https://www.hotstar.com").Retry(2, 5).End()
		if resp.StatusCode == 301 {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		//fmt.Println(body)
		u := resp.Header.Get("Location")
		if u == "" {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		t := strings.SplitN(u, "/", 4)
		if len(t) < 4 {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes, Region: t[3]}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.hotstar.com failed with code: %d", resp.StatusCode)}
}
