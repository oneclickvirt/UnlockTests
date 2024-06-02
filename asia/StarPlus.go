package asia

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// StarPlus
// www.starplus.com 双栈 且 get 请求
func StarPlus(request *gorequest.SuperAgent) model.Result {
	name := "Star+"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://www.starplus.com/"
	resp, body, errs := request.Get(url).
		Set("User-Agent", model.UA_Browser).
		Retry(2, 5).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Request.URL.String())
	if resp.StatusCode == 200 {
		if resp.StatusCode == 302 || resp.Header.Get("Location") == "https://www.preview.starplus.com/unavailable" {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		re := regexp.MustCompile(`Region:\s+([A-Za-z]{2})`)
		matches := re.FindStringSubmatch(body)
		if len(matches) >= 2 {
			loc := strings.ToLower(matches[1])
			if utils.GetRegion(loc, model.StarPlusSupportCountry) {
				anotherCheck := AnotherStarPlus()
				if anotherCheck.Err == nil && anotherCheck.Status == model.StatusYes {
					return model.Result{Name: name, Status: model.StatusYes, Region: loc}
				} else {
					anotherCheck.Info = "Website: " + model.StatusYes
					return anotherCheck
				}
			}
			return model.Result{Name: name, Status: model.StatusNo}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.starplus.com failed with code: %d", resp.StatusCode)}
}

// AnotherStarPlus
// StarPlus 的 另一个检测逻辑
func AnotherStarPlus() model.Result {
	name := "Star+"
	starcontent := "{\"query\":\"mutation registerDevice($input: RegisterDeviceInput!) " +
		"{\\n            registerDevice(registerDevice: $input) {\\n                grant " +
		"{\\n                    grantType\\n                    assertion\\n                " +
		"}\\n            }\\n        }\",\"variables\":{\"input\":{\"deviceFamily\":\"browser\"," +
		"\"applicationRuntime\":\"chrome\",\"deviceProfile\":\"windows\",\"deviceLanguage\":\"zh-CN\"," +
		"\"attributes\":{\"osDeviceIds\":[],\"manufacturer\":\"microsoft\",\"model\":null," +
		"\"operatingSystem\":\"windows\",\"operatingSystemVersion\":\"10.0\",\"browserName\":" +
		"\"chrome\",\"browserVersion\":\"96.0.4664\"}}}}"
	request := gorequest.New()
	resp, body, errs := request.Post("https://star.api.edge.bamgrid.com/graph/v1/device/graphql").
		Set("authorization", "c3RhciZicm93c2VyJjEuMC4w.COknIGCR7I6N0M5PGnlcdbESHGkNv7POwhFNL-_vIdg").
		SendString(starcontent).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	if resp.StatusCode >= 400 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	region := ""
	inSupportedLocation := false
	if d, ok := data["data"].(map[string]interface{}); ok {
		if r, ok := d["registerDevice"].(map[string]interface{}); ok {
			if g, ok := r["grant"].(map[string]interface{}); ok {
				if grantType, ok := g["grantType"].(string); ok && grantType == "jwt" {
					region = "UNKNOWN"
				}
			}
		}
	}
	isUnavailable := false
	if resp.StatusCode == 404 {
		isUnavailable = true
	}
	if region != "" && !isUnavailable && !inSupportedLocation {
		return model.Result{Name: name, Status: "CDN Relay Available"}
	} else if region != "" && isUnavailable {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if region != "" && inSupportedLocation {
		return model.Result{Name: name, Status: model.StatusYes, Region: region}
	} else if region == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.starplus.com another check failed with code: %d", resp.StatusCode)}
}
