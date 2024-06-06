package eu

import (
	"crypto/rand"
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"strings"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// MathsSpot
func MathsSpot(c *http.Client) model.Result {
	name := "Maths Spot"
	if c == nil {
		return model.Result{Name: name}
	}
	headers := map[string]string{
		"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language": "en-US,en;q=0.9",
		"User-Agent":      model.UA_Browser,
	}
	url := "https://mathsspot.com/"
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Get(url).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if len(body) > 0 && strings.Contains(body, "FailureServiceNotInRegion") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	region := utils.ReParse(body, `"countryCode"\s{0,}:\s{0,}"([^"]+)"`)
	nggFeVersion := utils.ReParse(body, `"NEXT_PUBLIC_FE_VERSION"\s{0,}:\s{0,}"([^"]+)"`)
	if region == "" || nggFeVersion == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	request2 := utils.Gorequest(c)
	headers2 := map[string]string{
		"accept":                "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":       "en-US,en;q=0.9",
		"User-Agent":            model.UA_Browser,
		"referer":               "https://mathsspot.com/",
		"sec-ch-ua":             model.UA_SecCHUA,
		"sec-fetch-dest":        "empty",
		"sec-fetch-mode":        "cors",
		"sec-fetch-site":        "same-origin",
		"x-ngg-skip-evar-check": "true",
		"x-ngg-fe-version":      nggFeVersion,
	}
	request2 = utils.SetGoRequestHeaders(request2, headers2)
	url2 := fmt.Sprintf("https://mathsspot.com/3/api/play/v1/startSession?appId=5349&uaId=ua-%s&uaSessionId=uasess-%s&feSessionId=fesess-%s&visitId=visitid-%s&initialOrientation=landscape&utmSource=NA&utmMedium=NA&utmCampaign=NA&deepLinkUrl=&accessCode=&ngReferrer=NA&pageReferrer=NA&ngEntryPoint=https%%3A%%2F%%2Fmathsspot.com%%2F&ntmSource=NA&customData=&appLaunchExtraData=&feSessionTags=nowgg&sdpType=&eVar=&isIframe=false&feDeviceType=desktop&feOsName=window&userSource=direct&visitSource=direct&userCampaign=NA&visitCampaign=NA", generateRandomString(21), generateRandomString(21), generateRandomString(21), generateRandomString(21))
	resp2, body2, errs2 := request2.Get(url2).End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	//fmt.Println(body2)
	status := utils.ReParse(body2, `"status"\s{0,}:\s{0,}"([^"]+)"`)
	switch status {
	case "FailureServiceNotInRegion":
		return model.Result{Name: name, Status: model.StatusNo}
	case "FailureProxyUserLimitExceeded":
		return model.Result{Name: name, Status: model.StatusNo, Info: "Proxy/VPN Detected"}
	case "Success":
		return model.Result{Name: name, Status: model.StatusYes, Region: region}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected, Info: status}
	}
}
