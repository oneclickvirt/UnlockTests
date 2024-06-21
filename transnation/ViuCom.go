package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// ViuCom
// www.viu.com 仅 ipv4 且 get 请求
func ViuCom(c *http.Client) model.Result {
	name := "Viu.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.viu.com"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	//b, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	//}
	//body := string(b)
	location := fmt.Sprintf("%s", resp.Request.URL)
	if strings.Contains(location, "no-service") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if location != "" {
		regions := strings.Split(location, "/")
		if regions[len(regions)-1] == "no-service" {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(regions[len(regions)-1])}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.viu.com failed with code: %d", resp.StatusCode)}
}
