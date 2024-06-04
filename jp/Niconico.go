package jp

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

// Niconico
// www.nicovideo.jp 仅 ipv4 且 get 请求
func Niconico(c *http.Client) model.Result {
	name := "Niconico"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://www.nicovideo.jp/watch/so40278367" // 进击的巨人
	//url2 := "https://www.nicovideo.jp/watch/so23017073" // 假面骑士
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp1, body1, errs1 := request.Get(url1).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	if strings.Contains(body1, "同じ地域") || resp1.StatusCode == 403 || resp1.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	headers = map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "en-US,en;q=0.9",
		"sec-ch-ua":                 `"(Not(A:Brand";v="8", "Chromium";v="114", "Google Chrome";v="114")"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
	}
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get("https://live.nicovideo.jp/?cmnhd_ref=device=pc&site=nicolive&pos=header_servicelink&ref=WatchPage-Anchor").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// 查找第一个官方直播剧的ID
	splitted := strings.Split(body, "&quot;isOfficialChannelMemberFree&quot;:false")
	var liveID string
	for _, part := range splitted {
		if strings.Contains(part, "話") && !strings.Contains(part, "&quot;isOfficialChannelMemberFree&quot;:true") && !strings.Contains(part, "playerProgram") && !strings.Contains(part, "&quot;ON_AIR&quot;") {
			startIdx := strings.Index(part, "&quot;id&quot;:&quot;")
			if startIdx != -1 {
				startIdx += len("&quot;id&quot;:&quot;")
				endIdx := strings.Index(part[startIdx:], "&quot;")
				if endIdx != -1 {
					liveID = part[startIdx : startIdx+endIdx]
					break
				}
			}
		}
	}
	if liveID != "" {
		resp, body, errs = request.Get("https://live.nicovideo.jp/watch/" + liveID).End()
		if len(errs) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
		}
		defer resp.Body.Close()
		if strings.Contains(body, "notAllowedCountry") && resp1.StatusCode == 200 {
			return model.Result{Name: name, Status: model.StatusYes,
				Info: fmt.Sprintf("But Official Live is Unavailable. LiveID: %s", liveID)}
		}
		if resp1.StatusCode == 200 {
			return model.Result{Name: name, Status: model.StatusYes,
				Info: fmt.Sprintf("LiveID: %s", liveID)}
		}
	} else if resp1.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes,
			Info: "But Official Live is Unavailable"}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.nicovideo.jp failed with code: %d", resp.StatusCode)}
}
