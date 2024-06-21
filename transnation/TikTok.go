package transnation

import (
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// TikTok
// www.tiktok.com 仅 ipv4 且 get 请求
func TikTok(c *http.Client) model.Result {
	name := "TikTok"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.tiktok.com/explore"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	if resp.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if strings.Contains(body, "https://www.tiktok.com/hk/notfound") {
		return model.Result{Name: name, Status: model.StatusNo, Region: "hk"}
	}
	region := utils.ReParse(body, `"region":"(\w+)"`)
	if region != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
	} else {
		url = "https://www.tiktok.com/"
		resp2, err2 := client.R().Get(url)
		if err2 != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
		}
		defer resp2.Body.Close()
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
		}
		body = string(b)
		if resp.StatusCode != 200 {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if strings.Contains(body, "https://www.tiktok.com/hk/notfound") {
			return model.Result{Name: name, Status: model.StatusNo, Region: "hk"}
		}
		region = utils.ReParse(body, `"region":"(\w+)"`)
		if region != "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region)}
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
	// return model.Result{Name: name, Status: model.StatusUnexpected,
	// 	Err: fmt.Errorf("www.tiktok.com can not find region with resp code: %d", resp.StatusCode)}
}
