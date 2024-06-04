package jp

import (
	"encoding/json"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Abema
// api.abema.io 仅 ipv4 且 get 请求
func Abema(c *http.Client) model.Result {
	name := "Abema.TV"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.abema.io/v1/ip/check?device=android"
	headers := map[string]string{
		"User-Agent": model.UA_Dalvik,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// fmt.Println(body)
	var abemaRes struct {
		Message        string `json:"message"`
		IsoCountryCode string `json:"isoCountryCode"`
	}
	if err := json.Unmarshal([]byte(body), &abemaRes); err != nil {
		if strings.Contains(body, "blocked_location") || strings.Contains(body, "anonymous_ip") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if abemaRes.IsoCountryCode == "JP" || strings.Contains(body, "JP") {
		return model.Result{Name: name, Status: model.StatusYes, Region: "JP"}
	}
	if abemaRes.Message == "blocked_location" || abemaRes.Message == "anonymous_ip" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusYes + " (Oversea Only)"}
}
