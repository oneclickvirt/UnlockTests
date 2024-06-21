package au

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
)

// Channel9
// login.nine.com.au 双栈 且 get 请求
func Channel9(c *http.Client) model.Result {
	name := "Channel 9"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://login.nine.com.au"
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
	//fmt.Println(body)
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 302 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get login.nine.com.au failed with code: %d", resp.StatusCode)}
}
