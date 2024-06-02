package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// WikipediaEditable
// zh.wikipedia.org 双栈 get 请求
func WikipediaEditable(request *gorequest.SuperAgent) model.Result {
	name := "Wikipedia Editability"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://zh.wikipedia.org/w/index.php?title=Wikipedia%3A%E6%B2%99%E7%9B%92&action=edit"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "Banned") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get zh.wikipedia.org failed with code: %d", resp.StatusCode)}
}
