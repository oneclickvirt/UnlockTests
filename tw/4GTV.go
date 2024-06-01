package tw

import (
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"time"
)

// Tw4gtv
// api2.4gtv.tv 仅 ipv4 且 post 请求
func Tw4gtv(request *gorequest.SuperAgent) model.Result {
	name := "4GTV.TV"
	url := "https://api2.4gtv.tv//Vod/GetVodUrl3"
	formData := `value=D33jXJ0JVFkBqV%2BZSi1mhPltbejAbPYbDnyI9hmfqjKaQwRQdj7ZKZRAdb16%2FRUrE8vGXLFfNKBLKJv%2BfDSiD%2BZJlUa5Msps2P4IWuTrUP1%2BCnS255YfRadf%2BKLUhIPj`
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Post(url).
		Timeout(15 * time.Second).
		Type("form").
		Send(formData).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var res struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Success {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if res.Success == false || resp.StatusCode == 403 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api2.4gtv.tv failed with code: %d", resp.StatusCode)}
}
