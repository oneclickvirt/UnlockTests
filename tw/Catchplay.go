package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Catchplay
// sunapi.catchplay.com 仅 ipv4 且 get 请求
// unauthorized 有问题
func Catchplay(c *http.Client) model.Result {
	name := "CatchPlay+"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://sunapi.catchplay.com/geo"
	headers := map[string]string{
		"authorization": "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		if strings.Contains(body, "is not allowed") && strings.Contains(body, "The location") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		fmt.Println(body)
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Code == "100016" {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if res.Code == "0" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sunapi.catchplay.com failed with code: %d", resp.StatusCode)}
}
