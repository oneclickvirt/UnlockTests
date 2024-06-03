package jp

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

type PlatformResponse struct {
	PlatformUID   string `json:"platform_uid"`
	PlatformToken string `json:"platform_token"`
}

type Episode struct {
	AccountID  string `json:"accountID"`
	PlayerID   string `json:"playerID"`
	VideoID    string `json:"videoID"`
	VideoRefID string `json:"videoRefID"`
}

func getEpisodeID(body string) string {
	if idx := strings.Index(body, `"newer-drama"`); idx != -1 {
		body = body[idx:]
		if idx = strings.Index(body, `"id"`); idx != -1 {
			body = body[idx:]
			parts := strings.Split(body, `"`)
			if len(parts) > 3 {
				return parts[3]
			}
		}
	}
	return ""
}

func getPolicyKey(body string) string {
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
// 双重检测逻辑
func TVer(request *gorequest.SuperAgent) model.Result {
	firstCheck := FirstTVer(request)
	//if firstCheck.Status == model.StatusNetworkErr {
	//	secondCheck := AnotherTVer()
	//	if secondCheck.Status == model.StatusNetworkErr {
	//		return firstCheck
	//	} else {
	//		return secondCheck
	//	}
	//} else {
	//	return firstCheck
	//}
	return firstCheck
}

// FirstTVer
// 主要的检测逻辑
func FirstTVer(request *gorequest.SuperAgent) model.Result {
	name := "TVer"
	if request == nil {
		return model.Result{Name: name}
	}
	// 创建平台用户
	res, body, errs := request.Post("https://platform-api.tver.jp/v2/api/platform_users/browser/create").
		Set("content-type", "application/x-www-form-urlencoded").
		Set("origin", "https://s.tver.jp").
		Set("referer", "https://s.tver.jp/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("user-agent", model.UA_Browser).
		Send("device_type=pc").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
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

	// 获取当前播放的剧集
	url := fmt.Sprintf("https://platform-api.tver.jp/service/api/v1/callHome?"+
		"platform_uid=%s&platform_token=%s&require_data=mylist%%2Cresume%%2Clater", platformResp.PlatformUID,
		platformResp.PlatformToken)
	res, body, errs = request.Get(url).
		Set("origin", "https://tver.jp").
		Set("referer", "https://tver.jp/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("x-tver-platform-type", "web").
		Set("user-agent", model.UA_Browser).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("2. get platform-api.tver.jp failed with code: %d", res.StatusCode)}
	}
	episodeId := getEpisodeID(body)
	if episodeId == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (No Episode ID)")}
	}

	// 获取剧集的信息
	url = fmt.Sprintf("https://statics.tver.jp/content/episode/%s.json", episodeId)
	res, body, errs = request.Get(url).
		Set("origin", "https://tver.jp").
		Set("referer", "https://tver.jp/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("user-agent", model.UA_Browser).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("get platform-api.tver.jp failed with code: %d", res.StatusCode)}
	}
	var episode Episode
	if err := json.Unmarshal([]byte(body), &episode); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (Parsing JSON)")}
	}

	// 获取 Brightcove 播放器信息
	url = fmt.Sprintf("https://players.brightcove.net/%s/%s_default/index.min.js", episode.AccountID,
		episode.PlayerID)
	res, body, errs = request.Get(url).
		Set("Referer", "https://tver.jp/").
		Set("Sec-Fetch-Dest", "script").
		Set("Sec-Fetch-Mode", "no-cors").
		Set("Sec-Fetch-Site", "cross-site").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("user-agent", model.UA_Browser).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.Result{Name: name, Status: model.StatusNetworkErr,
			Err: fmt.Errorf("get platform-api.tver.jp failed with code: %d", res.StatusCode)}
	}
	policyKey := getPolicyKey(body)
	if policyKey == "" {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("failed (No policyKey)")}
	}

	// 最终测试
	var finalURL string
	if episode.VideoRefID == "" {
		deliveryConfigId := getDeliveryConfigID(body)
		if deliveryConfigId != "" {
			finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/%s?config_id=%s",
				episode.AccountID, episode.VideoID, deliveryConfigId)
		} else {
			finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/ref%%3A%s",
				episode.AccountID, episode.VideoRefID)
		}
	} else {
		finalURL = fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/ref%%3A%s",
			episode.AccountID, episode.VideoRefID)
	}

	// 构建请求
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("accept", fmt.Sprintf("application/json;pk=%s", policyKey))
	client.Headers.Set("origin", "https://tver.jp")
	client.Headers.Set("referer", "https://tver.jp/")
	client.Headers.Set("sec-ch-ua", model.UA_SecCHUA)
	client.Headers.Set("sec-ch-ua-mobile", "?0")
	client.Headers.Set("sec-ch-ua-platform", "Windows")
	client.Headers.Set("sec-fetch-dest", "empty")
	client.Headers.Set("sec-fetch-mode", "cors")
	client.Headers.Set("sec-fetch-site", "cross-site")
	client.Headers.Set("user-agent", model.UA_Browser)
	resp, err := client.R().Get(finalURL)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
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
	if err := json.Unmarshal(b, &res1); err != nil {
		if err := json.Unmarshal(b, &res2); err != nil {
			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res2[0].ErrorSubcode == "CLIENT_GEO" {
			return model.Result{Name: name, Status: model.StatusNo, Region: res2[0].ClientGeo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res1.AccountId != "0" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "jp"}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get edge.api.brightcove.com failed with code: %d", resp.StatusCode)}
}

//// AnotherTVer
//func AnotherTVer() model.Result {
//	client := req.DefaultClient()
//	client.ImpersonateChrome()
//	client.Headers.Set("User-Agent", model.UA_Browser)
//	client.Headers.Set("Accept", "application/json;pk=BCpkADawqM0_rzsjsYbC1k1wlJLU4HiAtfzjxdUmfvvLUQB-Ax6VA-p-9wOEZbCEm3u95qq2Y1CQQW1K9tPaMma9iAqUqhpISCmyXrgnlpx9soEmoVNuQpiyGsTpePGumWxSs1YoKziYB6Wz")
//	resp, err := client.R().
// SetRetryCount(2).
// SetRetryBackoffInterval(1*time.Second, 5*time.Second).
// SetRetryFixedInterval(2 * time.Second).
// Get("https://edge.api.brightcove.com/playback/v1/accounts/5102072605001/videos/ref%3Akaguyasama_01")
//	if err != nil {
//		return model.Result{Status: model.StatusNetworkErr, Err: err}
//	}
//	defer resp.Body.Close()
//	b, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return model.Result{Status: model.StatusNetworkErr, Err: err}
//	}
//	body := string(b)
//	var res1 struct {
//		ErrorSubcode string `json:"error_subcode"`
//		AccountId    string `json:"account_id"`
//	}
//	var res2 []struct {
//		ClientGeo    string `json:"client_geo"`
//		ErrorSubcode string `json:"error_subcode"`
//		ErrorCode    string `json:"error_code"`
//		Message      string `json:"message"`
//	}
//	fmt.Println(body)
//	if err := json.Unmarshal(b, &res1); err != nil {
//		if err := json.Unmarshal(b, &res2); err != nil {
//			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
//				return model.Result{
//					Status: model.StatusNo,
//				}
//			}
//			return model.Result{Status: model.StatusErr, Err: err}
//		}
//		if res2[0].ErrorSubcode == "CLIENT_GEO" {
//			return model.Result{Status: model.StatusNo, Region: res2[0].ClientGeo}
//		}
//		return model.Result{Status: model.StatusErr, Err: err}
//	}
//	if res1.AccountId != "0" {
//		return model.Result{Status: model.StatusYes, Region: "jp"}
//	}
//	return model.Result{Status: model.StatusUnexpected,
//		Err: fmt.Errorf("get edge.api.brightcove.com failed with code: %d", resp.StatusCode)}
//}
