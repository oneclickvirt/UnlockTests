package ca

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)


// CBCGem
// www.cbc.ca 仅 ipv4 且 get 请求
func CBCGem(request *gorequest.SuperAgent) model.Result {
	name := "CBC Gem"
	url := "https://www.cbc.ca/g/stats/js/cbc-stats-top.js"
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, `country":"CA"`) {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	if resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.cbc.ca failed with code: %d", resp.StatusCode)}

}
