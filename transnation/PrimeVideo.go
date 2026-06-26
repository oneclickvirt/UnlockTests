package transnation

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

func PrimeVideo(c *http.Client) model.Result {
	name := "Amazon Prime Video"
	hostname := "www.primevideo.com"
	if c == nil {
		return model.Result{Name: name}
	}

	client := utils.Req(c)
	resp, err := client.R().Get("https://www.primevideo.com")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		if loc := resp.Header.Get("Location"); loc != "" {
			if strings.HasPrefix(loc, "/") {
				loc = "https://www.primevideo.com" + loc
			}
			if resp2, err := client.R().Get(loc); err == nil {
				defer resp2.Body.Close()
				if b2, readErr := io.ReadAll(resp2.Body); readErr == nil {
					body = string(b2)
				}
			}
		}
	}

	territoryRe := regexp.MustCompile(`"currentTerritory"\s*:\s*"([A-Za-z]{2})"`)
	if !territoryRe.MatchString(body) {
		storefrontRe := regexp.MustCompile(`(https://www\.amazon\.[a-z.]+/[^"'\s>]+)`)
		for _, match := range storefrontRe.FindAllStringSubmatch(body, -1) {
			urlStr := strings.ReplaceAll(match[1], "&amp;", "&")
			urlStr = strings.ReplaceAll(urlStr, `\u0026`, "&")
			if !strings.Contains(urlStr, "storefront") {
				continue
			}
			resp2, err := client.R().Get(urlStr)
			if err != nil {
				continue
			}
			b2, readErr := io.ReadAll(resp2.Body)
			resp2.Body.Close()
			if readErr == nil {
				body = string(b2)
				break
			}
		}
	}

	if strings.Contains(body, "api-services-support@amazon.com") {
		return model.Result{Name: name, Status: model.StatusNo}
	}

	if match := territoryRe.FindStringSubmatch(body); len(match) > 1 {
		location := strings.ToLower(match[1])
		if location != "cn" && location != "cu" && location != "ir" && location != "kp" && location != "sy" {
			result1, result2, result3 := utils.CheckDNS(hostname)
			unlockType := utils.GetUnlockType(result1, result2, result3)
			return model.Result{
				Name:       name,
				Status:     model.StatusYes,
				Region:     location,
				UnlockType: unlockType,
			}
		}
		return model.Result{Name: name, Status: model.StatusNo, Region: location}
	}

	return model.Result{Name: name, Status: model.StatusNo}
}
