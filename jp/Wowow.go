package jp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
	"time"
)

func getFirstLink(jsonStr string) string {
	type Drama struct {
		Link string `json:"link"`
	}
	var dramas []Drama
	err := json.Unmarshal([]byte(jsonStr), &dramas)
	if err != nil || len(dramas) == 0 {
		return ""
	}
	return dramas[0].Link
}

func getWodUrl(htmlStr string) string {
	lines := strings.Split(htmlStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "https://wod.wowow.co.jp/content/") {
			tempList := strings.Split(line, "https://wod.wowow.co.jp/content/")
			if len(tempList) >= 2 {
				tpList := strings.Split(tempList[1], "\"")
				if len(tpList) >= 2 {
					return "https://wod.wowow.co.jp/content/" + tpList[0]
				} else {
					return ""
				}
			} else {
				return ""
			}
		}
	}
	return ""
}

func getProgramUrl(htmlStr string) string {
	lines := strings.Split(htmlStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "https://wod.wowow.co.jp/program/") {
			tempList := strings.Split(line, ":")
			if len(tempList) >= 2 {
				tpList := strings.Split(tempList[len(tempList)-1], "\"")
				if len(tpList) >= 2 {
					for _, l := range tpList {
						if strings.Contains(l, "//wod.wowow.co.jp/program/") {
							return "https:" + l
						}
					}
				} else {
					return ""
				}
			} else {
				return ""
			}
		}
	}
	return ""
}

func getMetaId(htmlStr string) string {
	lines := strings.Split(htmlStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "https://wod.wowow.co.jp/watch/") {
			tempList := strings.Split(line, "https://wod.wowow.co.jp/watch/")
			if len(tempList) >= 2 {
				tpList := strings.Split(tempList[1], "\"")
				if len(tpList) >= 2 {
					return tpList[0]
				} else {
					return ""
				}
			} else {
				return ""
			}
		}
	}
	return ""
}

// Wowow
// www.wowow.co.jp 仅 ipv4 且 get 请求
func Wowow(c *http.Client) model.Result {
	name := "WOWOW"
	if c == nil {
		return model.Result{Name: name}
	}
	// 获取当前时间戳
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	// 第一次请求：获取原创剧集列表
	url := fmt.Sprintf("https://www.wowow.co.jp/drama/original/json/lineup.json?_=%d", timestamp)
	headers := map[string]string{
		"Accept":             "application/json, text/javascript, */*; q=0.01",
		"Referer":            "https://www.wowow.co.jp/drama/original/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"X-Requested-With":   "XMLHttpRequest",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"User-Agent":         model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr + " 1", Err: errs[0]}
	}
	defer resp.Body.Close()
	// 获取第一个剧集的链接
	playUrl := getFirstLink(body)
	if playUrl == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get play URL")}
	}

	// 第二次请求：获取真实链接
	headers2 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request2 := utils.Gorequest(c)
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	resp, body, errs = request2.Get(playUrl).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr + " 2", Err: errs[0]}
	}

	// 获取真实链接
	wodUrl := getWodUrl(body)
	if wodUrl == "" {
		programUrl := getProgramUrl(body)
		// 第二次请求的二次请求：获取真实链接
		headers3 := map[string]string{
			"User-Agent": model.UA_Browser,
		}
		request3 := utils.Gorequest(c)
		request3 = utils.SetGoRequestHeaders(request3, headers3)
		resp, body, errs = request3.Get(programUrl).End()
		if len(errs) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr + " 2-2", Err: errs[0]}
		}
		tempList := strings.Split(body, "\"refId\":\"")
		if len(tempList) >= 2 {
			for _, l := range tempList {
				if strings.Contains(l, "media_meta") {
					tpList := strings.Split(l, "\"")
					if len(tpList) >= 2 {
						wodUrl = "https://wod.wowow.co.jp/content/" + tpList[0]
						break
					}
				}
			}
		}
		if wodUrl == "" {
			return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get WOD URL")}
		}
	}

	// 第三次请求：获取 meta_id
	resp, body, errs = request.Get(wodUrl).
		Set("User-Agent", model.UA_Browser).
		Timeout(10 * time.Second).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr + " 3", Err: errs[0]}
	}
	metaId := getMetaId(body)
	if metaId == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get meta ID")}
	}

	// 生成 vUid
	hash := md5.Sum([]byte(fmt.Sprintf("%d", timestamp)))
	vUid := hex.EncodeToString(hash[:])
	// 最终测试请求
	authUrl := "https://mapi.wowow.co.jp/api/v1/playback/auth"
	data := fmt.Sprintf(`{"meta_id":"%s","vuid":"%s","device_code":1,"app_id":1,"ua":"%s"}`, metaId, vUid, model.UA_Browser)
	headers4 := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"content-type":       "application/json;charset=UTF-8",
		"origin":             "https://wod.wowow.co.jp",
		"referer":            "https://wod.wowow.co.jp/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"x-requested-with":   "XMLHttpRequest",
		"User-Agent":         model.UA_Browser,
	}
	request4 := utils.Gorequest(c)
	request4 = utils.SetGoRequestHeaders(request4, headers4)
	resp, body, errs = request.Post(authUrl).Send(data).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr + " 4", Err: errs[0]}
	}
	//fmt.Println(body)
	// {"error":{"message":"サポート外ネットワークからの接続です。日本国外からの接続、VPN・プロキシ経由の接続等ではご利用いただけません。","code":2055,"type":"Forbidden",
	if strings.Contains(body, "VPN") || strings.Contains(body, "Forbidden") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body, "playback_session_id") {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{
		Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mapi.wowow.co.jp failed with code: %d", resp.StatusCode)}
}
