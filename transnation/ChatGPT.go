package transnation

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// OpenAI
// api.openai.com 仅 ipv4 且 get 请求
func OpenAI(request *gorequest.SuperAgent) model.Result {
	name := "ChatGPT"
	if request == nil {
		return model.Result{Name: name}
	}
	url1 := "https://api.openai.com/compliance/cookie_requirements"
	resp1, body1, errs1 := request.Get(url1).
		Set("User-Agent", model.UA_Browser).
		Set("authority", "api.openai.com").
		Set("accept", "*/*").
		Set("accept-language", "zh-CN,zh;q=0.9").
		Set("authorization", "Bearer null").
		Set("content-type", "application/json").
		Set("origin", "https://platform.openai.com").
		Set("referer", "https://platform.openai.com/").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "empty").
		Set("sec-fetch-mode", "cors").
		Set("sec-fetch-site", "same-site").
		End()

	url2 := "https://ios.chat.openai.com/"
	resp2, body2, errs2 := request.Get(url2).
		Set("User-Agent", model.UA_Browser).
		Set("authority", "ios.chat.openai.com").
		Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Set("accept-language", "zh-CN,zh;q=0.9").
		Set("sec-ch-ua", model.UA_SecCHUA).
		Set("sec-ch-ua-mobile", "?0").
		Set("sec-ch-ua-platform", "Windows").
		Set("sec-fetch-dest", "document").
		Set("sec-fetch-mode", "navigate").
		Set("sec-fetch-site", "none").
		Set("sec-fetch-user", "?1").
		Set("upgrade-insecure-requests", "1").
		End()

	url3 := "https://chat.openai.com/cdn-cgi/trace"
	resp3, body3, errs3 := request.Get(url3).End()

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
				return model.Result{Name: name, Status: "429 Rate limit", Region: loc}
			}
		}
		return model.Result{Name: name, Status: "429 Rate limit"}
	}
	if !VPN && !unsupportedCountry && reqStatus1 && reqStatus2 && reqStatus3 {
		if location != "" {
			loc := strings.ToLower(location)
			exit := utils.GetRegion(loc, model.GptSupportCountry)
			if exit {
				return model.Result{Name: name, Status: model.StatusYes, Region: loc}
			} else {
				return model.Result{Name: name, Status: "Yes but cdn-cgi recognizes it unsupported", Region: location}
			}
		} else {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	} else if !unsupportedCountry && VPN && reqStatus1 {
		return model.Result{Name: name, Status: "Only Available with Web Browser"}
	} else if unsupportedCountry && !VPN && reqStatus2 {
		return model.Result{Name: name, Status: "Only Available with Mobile APP"}
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
