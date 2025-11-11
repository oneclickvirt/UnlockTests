package tw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// FridayVideo
// video.friday.tw 仅 ipv4 且 get 请求
func FridayVideo(c *http.Client) model.Result {
	name := "FridayVideo"
	hostname := "video.friday.tw"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://video.friday.tw/api2/streaming/get?streamingId=122581&streamingType=2&contentType=4&contentId=1&clientId="
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	var res struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return model.Result{Status: model.StatusErr, Err: err}
	}
	if res.Code != "null" {
		switch res.Code {
		case "1006":
			return model.Result{Name: name, Status: model.StatusNo}
		case "0000":
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
		default:
			return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected code: %s", res.Code)}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected response body: %s", body)}
}
