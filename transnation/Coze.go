package transnation

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func Coze(c *http.Client) model.Result {
	name := "Coze"
	hostname := "www.coze.com"
	if c == nil {
		return model.Result{Name: name}
	}
	resp, err := utils.Req(c).R().Get("https://www.coze.com/api/developer/get_login_info")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	var response struct {
		Code int `json:"code"`
		Data struct {
			IsForbiddenRegion bool   `json:"IsForbiddenRegion"`
			CountryCode       string `json:"CountryCode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &response); err != nil {
		body := string(b)
		if strings.Contains(body, "Your region is not supported") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if response.Data.IsForbiddenRegion {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if response.Data.CountryCode != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(response.Data.CountryCode)}
	}
	if response.Code == 0 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected}
}
