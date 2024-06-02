package nl

import (
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// ZIETCDN
// nlziet.nl 仅 ipv4 且 get 请求
// 直接通过CDN判断地区
func ZIETCDN() model.Result {
	name := "NLZIET"
	request := gorequest.New()
	resp, body, errs := request.Get("https://nlziet.nl/cdn-cgi/trace").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	tempList := strings.Split(body, "\n")
	var location string
	for _, line := range tempList {
		if strings.HasPrefix(line, "loc=") {
			location = strings.ReplaceAll(line, "loc=", "")
		}
	}
	loc := strings.ToLower(location)
	exit := utils.GetRegion(loc, model.NLZIETSupportCountry)
	if exit {
		return model.Result{Name: name, Status: model.StatusYes, Region: loc}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}

// ZIET
// nlziet.nl 仅 ipv4 且 get 请求 cookie 有效期非常短
func ZIET(request *gorequest.SuperAgent) model.Result {
	name := "NLZIET"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api.nlziet.nl/v7/stream/handshake/Widevine/Dash/VOD/rzIL9rb-TkSn-ek_wBmvaw?playerName=BitmovinWeb"
	resp, body, errs := request.Get(url).
		Set("User-Agent", model.UA_Browser).
		Set("Accept", "application/json, text/plain, */*").
		Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6").
		Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IkM4M0YzQUFGOTRCOTM0ODA2NkQwRjZDRTNEODhGQkREIiwidHlwIjoiYXQrand0In0.eyJuYmYiOjE3MTIxMjY0NTMsImV4cCI6MTcxMjE1NTI0OCwiaXNzIjoiaHR0cHM6Ly9pZC5ubHppZXQubmwiLCJhdWQiOiJhcGkiLCJjbGllbnRfaWQiOiJ0cmlwbGUtd2ViIiwic3ViIjoiMDAzMTZiNGEtMDAwMC0wMDAwLWNhZmUtZjFkZTA1ZGVlZmVlIiwiYXV0aF90aW1lIjoxNzEyMTI2NDUzLCJpZHAiOiJsb2NhbCIsImVtYWlsIjoibXVsdGkuZG5zMUBvdXRsb29rLmNvbSIsInVzZXJJZCI6IjMyMzg3MzAiLCJjdXN0b21lcklkIjoiMCIsImRldmljZUlkZW50aWZpZXIiOiJJZGVudGl6aWV0LTI0NWJiNmYzLWM2ZjktNDNjZS05ODhmLTgxNDc2OTcwM2E5OCIsImV4dGVybmFsVXNlcklkIjoiZTM1ZjdkMzktMjQ0ZC00ZTkzLWFkOTItNGFjYzVjNGY0NGNlIiwicHJvZmlsZUlkIjoiMjdDMzM3RjktOTRDRS00NjBDLTlBNjktMTlDNjlCRTYwQUIzIiwicHJvZmlsZUNvbG9yIjoiRkY0MjdDIiwicHJvZmlsZVR5cGUiOiJBZHVsdCIsIm5hbWUiOiJTdHJlYW1pbmciLCJqdGkiOiI4Q0M1QzYzNkJGRjg3MEE2REJBOERBNUMwQTk0RUZDRiIsImlhdCI6MTcxMjEyNjQ1Mywic2NvcGUiOlsiYXBpIiwib3BlbmlkIl0sImFtciI6WyJwcm9maWxlIiwicHdkIl19.bk-ziFPJM00bpE7TcgPmIYFFx-2Q5N3BkUzEvQ_dDMK9O1F9f7DEe-Qzmnb5ym7ChlnXwrCV3QyOOA24hu_gCrlNlD7-vI3XGZR-54zFD-F7cRDOoL-1-iO_10tmgwb5Io-svY0bn0EDYKeRxYYBi0w_3bFVFDM2CxxA6tWeBYIfN5rCSzBHd3RPPjYtqX-sogyh_5W_7KJ83GK5kpsywT3mz8q7Cs1mtKs9QA1-o01N0RvTxZAcfzsHg3-qGgLnvaAuZ_XqRK9kLWqJWeJTWKWtUI6OlPex22sY3keKFpfZnUtFv-BvkCM6tvbIlMZAClk3lhI8rMFAWDpUcbcS3w").
		Set("nlziet-appname", "WebApp").
		Set("nlziet-appversion", "5.43.24").
		Set("Origin", "https://app.nlziet.nl").
		Set("Referer", "https//app.nlziet.nl/").
		Set("Sec-Ch-UA", model.UA_SecCHUA).
		Set("Sec-Ch-UA-Mobile", "?0").
		Set("Sec-Ch-UA-Platform", "\"Windows\"").
		Set("Sec-Fetch-Dest", "empty").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site", "same-site").
		Timeout(10).End()
	if len(errs) > 0 {
		return ZIETCDN()
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		if strings.Contains(body, "CountryNotAllowed") {
			return model.Result{Name: name, Status: model.StatusNo}
		} else if strings.Contains(body, "streamSessionId") {
			return model.Result{Name: name, Status: model.StatusYes}
		} else {
			return ZIETCDN()
		}
	} else {
		return ZIETCDN()
	}
}
