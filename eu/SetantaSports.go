package eu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

const defaultSetantaAPIKey = "857a1e5d-e35e-4fdf-805b-a87b6f8364bf"

// SetantaSports
// dce-frontoffice.imggaming.com 仅 ipv4 且 get 请求
func SetantaSports(c *http.Client) model.Result {
	name := "Setanta Sports"
	hostname := "imggaming.com"
	if c == nil {
		return model.Result{Name: name}
	}
	apiKey := strings.TrimSpace(os.Getenv("UNLOCKTESTS_SETANTA_API_KEY"))
	if apiKey == "" {
		apiKey = defaultSetantaAPIKey
	}
	url := "https://dce-frontoffice.imggaming.com/api/v2/consent-prompt"
	headers := map[string]string{
		"User-Agent":      model.UA_Browser,
		"Realm":           "dce.adjara",
		"x-api-key":       apiKey,
		"Accept-Language": "en-US",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	var consentResponse struct {
		OutsideAllowedTerritories bool `json:"outsideAllowedTerritories"`
	}
	if err := json.Unmarshal(b, &consentResponse); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if consentResponse.OutsideAllowedTerritories {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if !consentResponse.OutsideAllowedTerritories {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{
		Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get dce-frontoffice.imggaming.com failed with code: %d", resp.StatusCode)}
}
