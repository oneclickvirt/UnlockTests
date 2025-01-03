package jp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
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
				}
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
				}
			}
		}
	}
	return ""
}

// 通过 program ID 获取详细信息
func getProgramDetails(c *http.Client, programID string) (string, error) {
	apiURL := "https://www.wowow.co.jp/API/new_prg/programdetail.php"
	data := fmt.Sprintf(`{"prg_cd": "%s", "mode": "19"}`, programID)
	headers := map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
		"Accept":       "application/json, text/plain, */*",
		"Origin":       "https://www.wowow.co.jp",
		"Referer":      "https://www.wowow.co.jp/",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().SetBody(data).Post(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from program details API")
	}
	var res struct {
		ArchiveURL string `json:"archive_url"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return "", fmt.Errorf("failed to parse program details response: %v", err)
	}
	return res.ArchiveURL, nil
}

// 尝试使用推荐列表方法获取 archiveURL
func tryRecommendListMethod(c *http.Client) (string, error) {
	url := "https://www.wowow.co.jp/assets/config/top_recommend_list.json"
	headers := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Referer":          "https://www.wowow.co.jp/",
		"X-Requested-With": "XMLHttpRequest",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res struct {
		DramaOriginal []interface{} `json:"drama_original"`
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return "", err
	}
	// 遍历所有剧集ID尝试获取archiveURL
	for _, id := range res.DramaOriginal {
		var programID string
		if str, ok := id.(string); ok {
			programID = str
		} else if num, ok := id.(float64); ok {
			programID = strconv.FormatFloat(num, 'f', 0, 64)
		} else {
			continue
		}
		archiveURL, err := getProgramDetails(c, programID)
		if err == nil && archiveURL != "" {
			return archiveURL, nil
		}
	}
	return "", fmt.Errorf("no valid archive URL found in recommend list")
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
				}
			}
		}
	}
	return ""
}

// Wowow
func Wowow(c *http.Client) model.Result {
	name := "WOWOW"
	hostname := "wowow.co.jp"
	if c == nil {
		return model.Result{Name: name}
	}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	var wodUrl string
	var firstMethodFailed bool
	// 第一种方法：通过原创剧集列表
	url := fmt.Sprintf("https://www.wowow.co.jp/drama/original/json/lineup.json?_=%d", timestamp)
	headers := map[string]string{
		"Accept":             "application/json, text/javascript, */*; q=0.01",
		"Referer":            "https://www.wowow.co.jp/drama/original/",
		"X-Requested-With":   "XMLHttpRequest",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"User-Agent":         model.UA_Browser,
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err == nil {
		b, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			playUrl := getFirstLink(string(b))
			if playUrl != "" {
				resp2, err := client.R().Get(playUrl)
				if err == nil {
					b2, err := io.ReadAll(resp2.Body)
					resp2.Body.Close()
					if err == nil {
						programID := getProgramUrl(string(b2))
						if programID != "" {
							archiveURL, err := getProgramDetails(c, programID)
							if err == nil {
								wodUrl = archiveURL
							}
						}
					}
				}
			}
		}
	}
	// 如果第一种方法失败，尝试第二种方法
	if wodUrl == "" {
		firstMethodFailed = true
		archiveURL, err := tryRecommendListMethod(c)
		if err == nil {
			wodUrl = archiveURL
		}
	}
	// 如果两种方法都失败了
	if wodUrl == "" {
		if firstMethodFailed {
			return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("both methods failed to get wod URL")}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get wod URL")}
	}
	// 获取 meta_id
	resp3, err := client.R().Get(wodUrl)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "3", Err: err}
	}
	b3, err := io.ReadAll(resp3.Body)
	resp3.Body.Close()
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "3", Err: fmt.Errorf("can not parse body")}
	}
	metaId := getMetaId(string(b3))
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
	resp4, body4, err := utils.PostJson(c, authUrl, data, headers4)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "4", Err: err}
	}
	if strings.Contains(body4, "VPN") || strings.Contains(body4, "Forbidden") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if strings.Contains(body4, "playback_session_id") {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{
		Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mapi.wowow.co.jp failed with code: %d", resp4.StatusCode)}
}
