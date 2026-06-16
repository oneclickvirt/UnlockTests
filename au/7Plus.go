package au

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

type sevenPlusMarketResponse struct {
	ID        int    `json:"_id"`
	Postcode  string `json:"postcode"`
	PlaceName string `json:"place_name"`
	IPAddress string `json:"ip_address"`
}

// Au7plus
// 7plus.com.au 通过官方市场定位接口判断是否为澳大利亚市场
func Au7plus(c *http.Client) model.Result {
	name := "7plus"
	siteHostname := "7plus.com.au"
	apiHostname := "market-cdn.swm.digital"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://market-cdn.swm.digital/v1/market/ip/?apikey=web")
	if err != nil {
		return utils.HandleNetworkError(c, apiHostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusUnexpected,
			Err: fmt.Errorf("get 7plus market failed with code: %d", resp.StatusCode)}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	var market sevenPlusMarketResponse
	if err := json.Unmarshal(body, &market); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	unlockType := ""
	if market.ID == 4 {
		result1, result2, result3 := utils.CheckDNS(siteHostname)
		unlockType = utils.GetUnlockType(result1, result2, result3)
	}
	return evaluateSevenPlusMarket(name, market, unlockType)
}

func evaluateSevenPlusMarket(name string, market sevenPlusMarketResponse, unlockType string) model.Result {
	info := market.PlaceName
	if info == "" && market.ID != 0 {
		info = fmt.Sprintf("market:%d", market.ID)
	}
	if market.ID != 4 {
		return model.Result{Name: name, Status: model.StatusNo, Info: info}
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: "au", Info: info, UnlockType: unlockType}
}
