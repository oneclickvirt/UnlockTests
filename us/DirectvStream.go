package us

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// ATTNOW - DirectvStream
// www.atttvnow.com 双栈 且 get 请求
func DirectvStream(c *http.Client) model.Result {
	name := "Directv Stream"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.atttvnow.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	// }
	// body := string(b)
	// fmt.Println(body)
	// fmt.Println(resp.Header)
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes}
}
