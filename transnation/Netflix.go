package transnation

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

func extractRegionFromPage(body string) string {
	re1 := regexp.MustCompile(`"country"\s*:\s*"([A-Z]{2})"`)
	if matches := re1.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	re2 := regexp.MustCompile(`"requestCountry"\s*:\s*\{\s*"id"\s*:\s*"([A-Z]{2})"`)
	if matches := re2.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	re3 := regexp.MustCompile(`"preferredLocale"\s*:\s*\{\s*"country"\s*:\s*"([A-Z]{2})"`)
	if matches := re3.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	re4 := regexp.MustCompile(`"geo"\s*:\s*\{[^}]*"country"\s*:\s*"([A-Z]{2})"`)
	if matches := re4.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	re5 := regexp.MustCompile(`data-country\s*=\s*"([A-Z]{2})"`)
	if matches := re5.FindStringSubmatch(body); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// NetflixCDN
// api.fast.com 双栈 get 请求
func NetflixCDN(c *http.Client) model.Result {
	name := "Netflix CDN"
	hostname := "fast.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api.fast.com/netflix/speedtest/v2?https=true&token=YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm&urlCount=5"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo, Info: "IP Banned By Netflix"}
	}
	type netflixCdnTarget struct {
		Name     string `json:"name"`
		Url      string `json:"url"`
		Location struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"location"`
	}
	var res struct {
		Targets []netflixCdnTarget `json:"targets"`
	}
	if err := json.Unmarshal(b, &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Targets[0].Location.Country != "" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{
			Name: name, Status: model.StatusYes,
			Region:     res.Targets[0].Location.Country,
			UnlockType: unlockType,
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.fast.com failed with code: %d", resp.StatusCode)}
}

// Netflix
// www.netflix.com 双栈 且 get 请求
func Netflix(c *http.Client) model.Result {
	name := "Netflix"
	hostname := "netflix.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://www.netflix.com/title/70143836" // 绝命毒师
	url2 := "https://www.netflix.com/title/81280792" // 乐高
	url3 := "https://www.netflix.com/title/80018499" // Test Patterns
	client1 := utils.Req(c)
	resp1, err1 := client1.R().Get(url1)
	if err1 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err1}
	}
	defer resp1.Body.Close()
	client2 := utils.Req(c)
	resp2, err2 := client2.R().Get(url2)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
	}
	defer resp2.Body.Close()
	if resp1.StatusCode == 404 && resp2.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusRestricted, Info: "Originals Only"}
	}
	if resp1.StatusCode == 403 && resp2.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if (resp1.StatusCode == 200 || resp1.StatusCode == 301) || (resp2.StatusCode == 200 || resp2.StatusCode == 301) {
		var bodyToCheck string
		var region string
		if resp1.StatusCode == 200 || resp1.StatusCode == 301 {
			b1, err := io.ReadAll(resp1.Body)
			if err != nil {
				return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
			}
			body1 := string(b1)
			hasVideo1 := strings.Contains(body1, `property="og:video"`)
			hasEpisodes1 := strings.Contains(body1, `data-uia="episodes"`)
			hasPlayableVideo1 := strings.Contains(body1, `playableVideo`)
			if hasVideo1 || hasEpisodes1 || hasPlayableVideo1 {
				bodyToCheck = body1
				region = extractRegionFromPage(body1)
			}
		}
		if bodyToCheck == "" && (resp2.StatusCode == 200 || resp2.StatusCode == 301) {
			b2, err := io.ReadAll(resp2.Body)
			if err != nil {
				return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
			}
			body2 := string(b2)
			hasVideo2 := strings.Contains(body2, `property="og:video"`)
			hasEpisodes2 := strings.Contains(body2, `data-uia="episodes"`)
			hasPlayableVideo2 := strings.Contains(body2, `playableVideo`)
			if hasVideo2 || hasEpisodes2 || hasPlayableVideo2 {
				bodyToCheck = body2
				region = extractRegionFromPage(body2)
			}
		}
		if bodyToCheck == "" {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if region == "" {
			client3 := utils.Req(c)
			resp3, err3 := client3.R().Get(url3)
			if err3 != nil {
				return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err3}
			}
			defer resp3.Body.Close()
			u := resp3.Header.Get("location")
			if u != "" {
				t := strings.SplitN(u, "/", 5)
				if len(t) >= 5 {
					region = strings.SplitN(t[3], "-", 2)[0]
				}
			}
			if region == "" {
				region = "Unknown"
			}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType, Region: strings.ToLower(region)}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.netflix.com failed with code: %d %d", resp1.StatusCode, resp2.StatusCode)}
}
