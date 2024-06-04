package transnation

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// GoogleSearch
// www.google.com 双栈 get 请求
func GoogleSearch(request *gorequest.SuperAgent) model.Result {
	name := "GoogleSearch"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.google.com/search?q=www.spiritysdx.top/"
	client := req.DefaultClient()
	client.ImpersonateChrome()
	resp, err := client.R().
		SetRetryCount(2).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second).
		Get(url)
	defer resp.Body.Close()
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	// fmt.Println(body)
	if strings.Contains(body, "unusual traffic from") || resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 && strings.Contains(body, "二叉树的博客") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.google.com failed with code: %d", resp.StatusCode)}
}
