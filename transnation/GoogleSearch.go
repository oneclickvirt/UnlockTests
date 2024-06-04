package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// GoogleSearch
// www.google.com 双栈 get 请求
func GoogleSearch(c *http.Client) model.Result {
	name := "GoogleSearch"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.google.com/search?q=www.spiritysdx.top/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
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