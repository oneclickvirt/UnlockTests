package jp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
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
func Wowow(request *gorequest.SuperAgent) model.Result {
	name := "WOWOW"
	// 获取当前时间戳
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	// 第一次请求：获取原创剧集列表
	url := fmt.Sprintf("https://www.wowow.co.jp/drama/original/json/lineup.json?_=%d", timestamp)
	resp, body, errs := request.Get(url).
		Set("Accept", "application/json, text/javascript, */*; q=0.01").
		Set("Referer", "https://www.wowow.co.jp/drama/original/").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "same-origin").
		Set("X-Requested-With", "XMLHttpRequest").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"Windows\"").
		Set("User-Agent", model.UA_Browser).
		Timeout(10 * time.Second).
		End()
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
	resp, body, errs = request.Get(playUrl).
		Set("User-Agent", model.UA_Browser).
		Timeout(10 * time.Second).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr + " 2", Err: errs[0]}
	}

	// 获取真实链接
	wodUrl := getWodUrl(body)
	if wodUrl == "" {
		programUrl := getProgramUrl(body)
		// 第二次请求的二次请求：获取真实链接
		resp, body, errs = request.Get(programUrl).
			Set("User-Agent", model.UA_Browser).
			Timeout(10 * time.Second).
			End()
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
	resp, body, errs = request.Post(authUrl).
		Set("accept", "application/json, text/plain, */*").
		Set("content-type", "application/json;charset=UTF-8").
		Set("origin", "https://wod.wowow.co.jp").
		Set("referer", "https://wod.wowow.co.jp/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "\"Windows\"").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		Set("x-requested-with", "XMLHttpRequest").
		Send(data).
		Set("User-Agent", model.UA_Browser).
		Timeout(10 * time.Second).
		End()
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
