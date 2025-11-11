package jp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

type PlatformResponse struct {
	Result struct {
		PlatformUID   string `json:"platform_uid"`
		PlatformToken string `json:"platform_token"`
	} `json:"result"`
}

type Episode struct {
	Video struct {
		AccountID  string `json:"accountID"`
		PlayerID   string `json:"playerID"`
		VideoID    string `json:"videoID"`
		VideoRefID string `json:"videoRefID"`
	} `json:"video"`
}

func getEpisodeID(body string) string {
	var homeResp struct {
		Result struct {
			Components []struct {
				ComponentID string `json:"componentID"`
				Contents    []struct {
					Content struct {
						EpisodeID string `json:"id"`
					} `json:"content"`
				} `json:"contents"`
			} `json:"components"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &homeResp); err == nil {
		for _, component := range homeResp.Result.Components {
			if component.ComponentID == "variety.catchup.recomend" && len(component.Contents) > 0 {
				for _, content := range component.Contents {
					id := content.Content.EpisodeID
					matched, _ := regexp.MatchString(`^[a-z0-9]{10}$`, id)
					if matched {
						return id
					}
				}
			}
		}
	}
	if idx := strings.Index(body, `"variety.catchup.recomend"`); idx != -1 {
		body = body[idx:]
		if idx = strings.Index(body, `"id"`); idx != -1 {
			body = body[idx:]
			parts := strings.Split(body, `"`)
			if len(parts) > 3 {
				id := parts[3]
				matched, _ := regexp.MatchString(`^[a-z0-9]{10}$`, id)
				if matched {
					return id
				}
			}
		}
	}
	return ""
}

func getPolicyKey(body string) string {
	re := regexp.MustCompile(`policyKey:"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	if idx := strings.Index(body, `policyKey:"`); idx != -1 {
		body = body[idx+len(`policyKey:"`):]
		parts := strings.Split(body, `"`)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

func getDeliveryConfigID(body string) string {
	re := regexp.MustCompile(`deliveryConfigId:"([^"]+)"`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return matches[1]
	}
	if idx := strings.Index(body, `deliveryConfigId:"`); idx != -1 {
		body = body[idx+len(`deliveryConfigId:"`):]
		parts := strings.Split(body, `"`)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return ""
}

// TVer
// edge.api.brightcove.com 仅 ipv4 且 get 请求
func TVer(c *http.Client) model.Result {
	firstCheck := FirstTVer(c)
	if firstCheck.Status == model.StatusNetworkErr || firstCheck.Status == model.StatusErr {
		secondCheck := AnotherTVer(c)
		if secondCheck.Status == model.StatusNetworkErr || secondCheck.Status == model.StatusErr {
			return firstCheck
		} else {
			return secondCheck
		}
	}
	return firstCheck
}

// FirstTVer
func FirstTVer(c *http.Client) model.Result {
	name := "TVer"
	hostname := "tver.jp"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"content-type":       "application/x-www-form-urlencoded",
		"origin":             "https://s.tver.jp",
		"referer":            "https://s.tver.jp/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         model.UA_Browser,
	}
	res, body, err := utils.PostJson(c, "https://platform-api.tver.jp/v2/api/platform_users/browser/create",
		"device_type=pc", headers)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("1. get platform-api.tver.jp failed with code: %d", res.StatusCode)}
	}
	var platformResp PlatformResponse
	if err := json.Unmarshal([]byte(body), &platformResp); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	url := fmt.Sprintf("https://platform-api.tver.jp/service/api/v1/callHome?"+
		"platform_uid=%s&platform_token=%s&require_data=mylist%%2Cresume%%2Clater", platformResp.Result.PlatformUID,
		platformResp.Result.PlatformToken)
	headers2 := map[string]string{
		"origin":               "https://tver.jp",
		"referer":              "https://tver.jp/",
		"sec-ch-ua":            model.UA_SecCHUA,
		"sec-ch-ua-mobile":     "?0",
		"sec-ch-ua-platform":   "Windows",
		"sec-fetch-dest":       "empty",
		"sec-fetch-mode":       "cors",
		"sec-fetch-site":       "same-site",
		"x-tver-platform-type": "web",
		"user-agent":           model.UA_Browser,
	}
	client2 := utils.Req(c)
	client2 = utils.SetReqHeaders(client2, headers2)
	resp2, err := client2.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp2.Body.Close()
	b, err := io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body2 := string(b)
	if resp2.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("2. get platform-api.tver.jp failed with code: %d", resp2.StatusCode)}
	}
	episodeId := getEpisodeID(body2)
	if episodeId == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (No Episode ID)")}
	}
	url = fmt.Sprintf("https://statics.tver.jp/content/episode/%s.json", episodeId)
	headers3 := map[string]string{
		"origin":             "https://tver.jp",
		"referer":            "https://tver.jp/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         model.UA_Browser,
	}
	client3 := utils.Req(c)
	client3 = utils.SetReqHeaders(client3, headers3)
	resp3, err := client3.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp3.Body.Close()
	b3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	if resp3.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("get statics.tver.jp failed with code: %d", resp3.StatusCode)}
	}
	var episode Episode
	if err := json.Unmarshal(b3, &episode); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (Parsing JSON)")}
	}
	url = fmt.Sprintf("https://players.brightcove.net/%s/%s_default/index.min.js", episode.Video.AccountID,
		episode.Video.PlayerID)
	headers4 := map[string]string{
		"Referer":            "https://tver.jp/",
		"Sec-Fetch-Dest":     "script",
		"Sec-Fetch-Mode":     "no-cors",
		"Sec-Fetch-Site":     "cross-site",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"user-agent":         model.UA_Browser,
	}
	client4 := utils.Req(c)
	client4 = utils.SetReqHeaders(client4, headers4)
	resp4, err := client4.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp4.Body.Close()
	b4, err := io.ReadAll(resp4.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body4 := string(b4)
	if resp4.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("get players.brightcove.net failed with code: %d", resp4.StatusCode)}
	}
	policyKey := getPolicyKey(body4)
	if policyKey == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (No policyKey)")}
	}
	var finalURL string
	deliveryConfigId := getDeliveryConfigID(body4)
	if episode.Video.VideoRefID == "" {
		if deliveryConfigId != "" {
			finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/%s?config_id=%s",
				episode.Video.AccountID, episode.Video.VideoID, deliveryConfigId)
		} else {
			finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/%s",
				episode.Video.AccountID, episode.Video.VideoID)
		}
	} else {
		finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/ref%%3A%s",
			episode.Video.AccountID, episode.Video.VideoRefID)
	}
	headers5 := map[string]string{
		"accept":             fmt.Sprintf("application/json;pk=%s", policyKey),
		"origin":             "https://tver.jp",
		"referer":            "https://tver.jp/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
		"user-agent":         model.UA_Browser,
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers5)
	resp, err := client.R().Get(finalURL)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	body = string(b)
	var res1 struct {
		ErrorSubcode string `json:"error_subcode"`
		AccountId    string `json:"account_id"`
	}
	var res2 []struct {
		ClientGeo    string `json:"client_geo"`
		ErrorSubcode string `json:"error_subcode"`
		ErrorCode    string `json:"error_code"`
		Message      string `json:"message"`
	}
	if err := json.Unmarshal(b, &res2); err != nil {
		if err := json.Unmarshal(b, &res1); err != nil {
			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res1.AccountId != "" && res1.AccountId != "0" {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: "jp"}
		}
	} else {
		if len(res2) > 0 && res2[0].ErrorSubcode == "CLIENT_GEO" {
			return model.Result{Name: name, Status: model.StatusNo, Region: res2[0].ClientGeo}
		}
		return model.Result{Name: name, Status: model.StatusErr}
	}

	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get edge.api.brightcove.com failed with code: %d", resp.StatusCode)}
}

// AnotherTVer
func AnotherTVer(c *http.Client) model.Result {
	name := "TVer"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://edge.api.brightcove.com/playback/v1/accounts/5102072605001/videos/ref%3Akaguyasama_01"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
		"Accept":     "application/json;pk=BCpkADawqM0_rzsjsYbC1k1wlJLU4HiAtfzjxdUmfvvLUQB-Ax6VA-p-9wOEZbCEm3u95qq2Y1CQQW1K9tPaMma9iAqUqhpISCmyXrgnlpx9soEmoVNuQpiyGsTpePGumWxSs1YoKziYB6Wz",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	var res1 struct {
		ErrorSubcode string `json:"error_subcode"`
		AccountId    string `json:"account_id"`
	}
	var res2 []struct {
		ClientGeo    string `json:"client_geo"`
		ErrorSubcode string `json:"error_subcode"`
		ErrorCode    string `json:"error_code"`
		Message      string `json:"message"`
	}
	if err := json.Unmarshal(b, &res2); err != nil {
		if err := json.Unmarshal(b, &res1); err != nil {
			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
				return model.Result{Status: model.StatusNo}
			}
			return model.Result{Status: model.StatusErr, Err: err}
		}
		if res1.AccountId != "" && res1.AccountId != "0" {
			return model.Result{Status: model.StatusYes, Region: "jp"}
		}
	} else {
		if len(res2) > 0 && res2[0].ErrorSubcode == "CLIENT_GEO" {
			return model.Result{Status: model.StatusNo, Region: res2[0].ClientGeo}
		}
		return model.Result{Status: model.StatusErr}
	}
	return model.Result{Status: model.StatusUnexpected,
		Err: fmt.Errorf("get edge.api.brightcove.com failed with code: %d", resp.StatusCode)}
}
