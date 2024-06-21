package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// ParamountPlus
// www.paramountplus.com 双栈 且 get 请求
func ParamountPlus(c *http.Client) model.Result {
	name := "Paramount+"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.paramountplus.com/"
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
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 && strings.Contains(body, "\"country_name_intl\":\"International\"") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.paramountplus.com failed with code: %d", resp.StatusCode)}
}
