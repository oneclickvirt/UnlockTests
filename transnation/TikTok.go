package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// TikTok
// www.tiktok.com 仅 ipv4 且 get 请求
func TikTok(c *http.Client) model.Result {
	name := "TikTok"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	url := "https://www.tiktok.com/explore"
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if strings.Contains(body, "https://www.tiktok.com/hk/notfound") {
		return model.Result{Name: name, Status: model.StatusNo, Region: "hk"}
	}
	region := utils.ReParse(body, `"region":"(\w+)"`)
	if region != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.tiktok.com failed with code: %d", resp.StatusCode)}
}
