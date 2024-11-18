package jp

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	hostname := "wowow.co.jp"
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
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "1-1", Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "1-1", Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)

	// 获取第一个剧集的链接
	playUrl := getFirstLink(body)
	if playUrl == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get play URL")}
	}

	// 第二次请求：获取真实链接
	resp2, err := client.R().Get(playUrl)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "2-1", Err: err}
	}
	defer resp2.Body.Close()
	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "2-1", Err: fmt.Errorf("can not parse body")}
	}
	body2 := string(b2)

	// 获取 program ID
	programID := getProgramUrl(body2)
	if programID == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get program ID")}
	}

	// 新增逻辑：通过 program ID 获取 archive URL
	archiveURL, err := getProgramDetails(c, programID)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to get archive URL: %v", err)}
	}

	// 使用 archive URL 获取 meta ID
	resp3, err := client.R().Get(archiveURL)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "3", Err: err}
	}
	defer resp3.Body.Close()
	b3, err := io.ReadAll(resp3.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Info: "3", Err: fmt.Errorf("can not parse body")}
	}
	body3 := string(b3)

	// 提取 meta_id
	metaId := getMetaId(body3)
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
