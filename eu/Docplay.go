package eu

import (
	"fmt"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Docplay
// - AU 、 New Zealand 、UK
// www.docplay.com 仅 ipv4 且 get 请求
func Docplay(request *gorequest.SuperAgent) model.Result {
	name := "Docplay"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.docplay.com/subscribe"
	request = request.Set("User-Agent", model.UA_Browser).Timeout(20 * time.Second)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "DocPlay hasn't launched in your part of the world yet.") ||
		resp.Request.URL.String() == "https://www.docplay.com/geoblocked" ||
		strings.Contains(resp.Header.Get("Set-Cookie"), "geoblocked=true") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.docplay.com failed with code: %d", resp.StatusCode)}
}
