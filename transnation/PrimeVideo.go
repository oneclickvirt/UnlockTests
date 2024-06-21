package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
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
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
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
