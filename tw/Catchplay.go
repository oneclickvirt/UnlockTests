package tw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Catchplay
// sunapi.catchplay.com 仅 ipv4 且 get 请求
func Catchplay(c *http.Client) model.Result {
	name := "CatchPlay+"
	hostname := "catchplay.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://sunapi.catchplay.com/geo"
	headers := map[string]string{
		"authorization": "Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	var res struct {
		Code string `json:"code"`
		Data struct {
			IsoCode string `json:"isoCode"`
		} `json:"data"`
	}
	// fmt.Println(body)
	if err = json.Unmarshal(b, &res); err != nil {
		if strings.Contains(body, "is not allowed") && strings.Contains(body, "The location") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Code == "100016" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	region := res.Data.IsoCode
	if res.Code == "0" || region != "" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get sunapi.catchplay.com failed with code: %d", resp.StatusCode)}
}
