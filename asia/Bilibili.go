package asia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Bilibili
// B站主体请求逻辑
func Bilibili(c *http.Client, name, url string) model.Result {
	hostname := "bilibili.com"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	if resp.StatusCode == http.StatusPreconditionFailed {
		return model.Result{Name: name, Status: model.StatusRestricted}
	}
	body := string(b)
	//fmt.Println(body)
	var res struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		if strings.Contains(body, "抱歉您所在地区不可观看") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if strings.Contains(body, "抱歉您所在地区不可观看") || strings.Contains(body, "The area is inaccessible") ||
		res.Code == 10004001 || res.Code == 10003003 || res.Code == -10403 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Code == 0 {
		result1, result2, result3 := utils.CheckDNS(hostname)
		//fmt.Println(result1, result2, result3)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.bilibili.com failed with code: %d", resp.StatusCode)}
}

func BilibiliAnime(c *http.Client) model.Result {
	name := "Bilibili Anime"
	hostname := "bilibili.com"
	if c == nil {
		return model.Result{Name: name}
	}
	client := utils.Req(c)
	resp, err := client.R().Get("https://api.bilibili.com/x/web-interface/zone")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	var res struct {
		Code int `json:"code"`
		Data struct {
			CountryCode int `json:"country_code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Code != 0 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	region := utils.CountryCodeToAlpha2(fmt.Sprint(res.Data.CountryCode))
	if region == "" {
		region = fmt.Sprint(res.Data.CountryCode)
	}

	testURL := ""
	switch strings.ToUpper(region) {
	case "HK", "MO":
		testURL = "https://api.bilibili.com/pgc/player/web/playurl?avid=473502608&cid=845838026&qn=0&type=&otype=json&ep_id=678506&fourk=1&fnver=0&fnval=16&module=bangumi"
	case "TW":
		testURL = "https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&module=bangumi"
	case "TH":
		testURL = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=10077726"
	case "ID":
		testURL = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11130043"
	case "VN":
		testURL = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11405745"
	case "MY", "SG", "PH", "BN", "KH", "LA", "MM", "TL":
		testURL = "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=347666"
	default:
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(region), UnlockType: unlockType}
	}

	testResult := Bilibili(c, name, testURL)
	testResult.Region = strings.ToLower(region)
	if testResult.Status == model.StatusNo {
		testResult.Status = model.StatusRestricted
	}
	return testResult
}
