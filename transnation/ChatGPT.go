package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)
// OpenAI
// api.openai.com 仅 ipv4 且 get 请求
func OpenAI(c *http.Client) model.Result {
	name := "ChatGPT"
	if c == nil {
		return model.Result{Name: name}
	}
	url1 := "https://api.openai.com/compliance/cookie_requirements"
	headers1 := map[string]string{
		"User-Agent":         model.UA_Browser,
		"authority":          "api.openai.com",
		"accept":             "*/*",
		"accept-language":    "zh-CN,zh;q=0.9",
		"authorization":      "Bearer null",
		"content-type":       "application/json",
		"origin":             "https://platform.openai.com",
		"referer":            "https://platform.openai.com/",
		"sec-ch-ua":          model.UA_SecCHUA,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
	}
	request1 := utils.Gorequest(c)
	request1 = utils.SetGoRequestHeaders(request1, headers1)
	resp1, body1, errs1 := request1.Get(url1).End()

	url2 := "https://ios.chat.openai.com/"
	headers2 := map[string]string{
		"User-Agent":                model.UA_Browser,
		"authority":                 "ios.chat.openai.com",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "zh-CN,zh;q=0.9",
		"sec-ch-ua":                 model.UA_SecCHUA,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "Windows",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}
	request2 := utils.Gorequest(c)
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	resp2, body2, errs2 := request2.Get(url2).End()

	url3 := "https://chat.openai.com/cdn-cgi/trace"
	headers3 := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	request3 := utils.Gorequest(c)
	request3 = utils.SetGoRequestHeaders(request3, headers3)
	resp3, body3, errs3 := request3.Get(url3).End()

	var reqStatus1, reqStatus2, reqStatus3 bool
	if len(errs1) > 0 {
		fmt.Println(errs1)
		reqStatus1 = false
	} else {
		reqStatus1 = true
		defer resp1.Body.Close()
	}
	if len(errs2) > 0 {
		fmt.Println(errs2)
		reqStatus2 = false
	} else {
		reqStatus2 = true
		defer resp2.Body.Close()
	}
	if len(errs3) > 0 {
		fmt.Println(errs3)
		reqStatus3 = false
	} else {
		reqStatus3 = true
		defer resp3.Body.Close()
	}
	unsupportedCountry := strings.Contains(body1, "unsupported_country")
	VPN := strings.Contains(body2, "VPN")
	tempList := strings.Split(body3, "\n")
	var location string
	if reqStatus3 {
		for _, line := range tempList {
			if strings.HasPrefix(line, "loc=") {
				location = strings.ReplaceAll(line, "loc=", "")
			}
		}
	}
	if (resp1 != nil && resp1.StatusCode == 429) || (resp2 != nil && resp2.StatusCode == 429) {
		if location != "" {
			loc := strings.ToLower(location)
			exit := utils.GetRegion(loc, model.GptSupportCountry)
			if exit {
				return model.Result{Name: name, Status: model.StatusNo, Info: "429 Rate limit", Region: loc}
			}
		}
		return model.Result{Name: name, Status: model.StatusNo, Info: "429 Rate limit"}
	}
	if !VPN && !unsupportedCountry && reqStatus1 && reqStatus2 && reqStatus3 {
		if location != "" {
			loc := strings.ToLower(location)
			exit := utils.GetRegion(loc, model.GptSupportCountry)
			if exit {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc}
			} else {
				return model.Result{Name: name, Status: model.StatusYes, Info: "but cdn-cgi not unsupported", Region: location}
			}
		} else {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	} else if !unsupportedCountry && VPN && reqStatus1 {
		return model.Result{Name: name, Status: model.StatusYes, Info: "Only Available with Web Browser"}
	} else if unsupportedCountry && !VPN && reqStatus2 {
		return model.Result{Name: name, Status: model.StatusYes, Info: "Only Available with Mobile APP"}
	} else if !reqStatus1 && VPN {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if VPN && unsupportedCountry {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if !reqStatus1 && !reqStatus2 && !reqStatus3 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
}
