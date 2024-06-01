package transnation

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// PrimeVideo
// www.primevideo.com 仅 ipv4 且 get 请求
func PrimeVideo(request *gorequest.SuperAgent) model.Result {
	name := "Amazon Prime Video"
	url := "https://www.primevideo.com/"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if i := strings.Index(body, `"currentTerritory":`); i != -1 {
		location := strings.ToLower(body[i+20 : i+22])
		if location != "cn" && location != "cu" && location != "ir" && location != "kp" && location != "sy" {
			return model.Result{
				Name: name, Status: model.StatusYes,
				Region: location,
			}
		}
		return model.Result{
			Name: name, Status: model.StatusNo,
			Region: location,
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
