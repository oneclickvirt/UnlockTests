package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// IQiYi
// www.iq.com 仅 ipv4 且 get 请求
func IQiYi(c *http.Client) model.Result {
	name := "IQiYi"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.iq.com"
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
