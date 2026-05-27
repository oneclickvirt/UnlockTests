package transnation

import (
	"encoding/json"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// DAZN
// startup.core.indazn.com 仅 ipv4 且 post 请求
func DAZN(c *http.Client) model.Result {
	name := "Dazn"
	hostname := "startup.core.indazn.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://www.dazn.com/" // check if 403 first
	client := utils.Req(c)
	resp, err := client.R().Get(url1)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	url2 := "https://startup.core.indazn.com/v1/main/web?Platform=web&LandingPageKey=generic&Brand=dazn" // check region
	resp2, err := client.R().Get(url2)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp2.Body.Close()
	var daznRes struct {
		Region struct {
			IsAllowed             bool   `json:"isAllowed"`
			DisallowedReason      string `json:"disallowedReason"`
			GeolocatedCountry     string `json:"GeolocatedCountry"`
			GeolocatedCountryName string `json:"GeolocatedCountryName"`
		} `json:"Region"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&daznRes); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if daznRes.Region.IsAllowed {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{
			Name: name, Status: model.StatusYes,
			Region:     daznRes.Region.GeolocatedCountry,
			UnlockType: unlockType,
		}
	}
	return model.Result{
		Name: name, Status: model.StatusNo, Info: daznRes.Region.DisallowedReason,
	}
}
