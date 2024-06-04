package transnation

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// PrimeVideo
// www.primevideo.com 仅 ipv4 且 get 请求
func PrimeVideo(c *http.Client) model.Result {
	name := "Amazon Prime Video"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.primevideo.com/"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
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
