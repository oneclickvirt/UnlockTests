package au

import (
	"fmt"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// KayoSports
// kayosports.com.au 实际使用 cf 检测，非澳洲请求将一直超时
func KayoSports(request *gorequest.SuperAgent) model.Result {
	name := "Kayo Sports"
	url := "https://kayosports.com.au"
	request = request.Set("User-Agent", model.UA_Browser).Timeout(20 * time.Second)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		// if strings.Contains(errs[0].Error(), "i/o timeout") {
		return model.Result{Name: name, Status: model.StatusNo}
		// }
		// return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(resp.Header.Get("Set-Cookie"), "geoblocked=true") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get kayosports.com.au failed with code: %d", resp.StatusCode)}
}
