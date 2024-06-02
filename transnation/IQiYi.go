package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// IQiYi
// www.iq.com 仅 ipv4 且 get 请求
func IQiYi(request *gorequest.SuperAgent) model.Result {
	name := "IQiYi"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.iq.com"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	tp := resp.Header.Get("x-custom-client-ip")
	if tp == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	var region string
	tpList := strings.Split(tp, ":")
	if len(tpList) >= 2 {
		region = tpList[len(tpList)-1]
		if region == "ntw" {
			region = "tw"
		}
	}
	if region != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: region}
	} else {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get www.iq.com failed with head: %s", tp)}
	}
}
