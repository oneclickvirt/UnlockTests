package us

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Crackle
// prod-api.crackle.com 双栈 get 请求
func Crackle(c *http.Client) model.Result {
	name := "Crackle"
	hostname := "crackle.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://prod-api.crackle.com/appconfig"
	headers := map[string]string{
		"User-Agent":           model.UA_Browser,
		"Accept-Language":      "en-US,en;q=0.9",
		"Content-Type":         "application/json",
		"Origin":               "https://www.crackle.com",
		"Referer":              "https://www.crackle.com/",
		"Sec-Fetch-Dest":       "empty",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Site":       "same-site",
		"sec-ch-ua":            "${UA_SEC_CH_UA}",
		"sec-ch-ua-mobile":     "?0",
		"sec-ch-ua-platform":   "\"Windows\"",
		"x-crackle-apiversion": "v2.0.0",
		"x-crackle-brand":      "crackle",
		"x-crackle-platform":   "5FE67CCA-069A-42C6-A20F-4B47A8054D46",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	// }
	// body := string(b)
	// fmt.Println(body)
	if resp.Header.Get("x-crackle-region") == "US" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: "us", UnlockType: unlockType}
	}
	if resp.Header.Get("x-crackle-region") != "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get prod-api.crackle.com failed with code: %d", resp.StatusCode)}
}
