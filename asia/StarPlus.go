package asia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

// StarPlus
// www.starplus.com 双栈 且 get 请求
func StarPlus(c *http.Client) model.Result {
	name := "Star+"
	hostname := "starplus.com"
	if c == nil {
		return model.Result{Name: name}
	}
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	url := "https://www.starplus.com/"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("StarPlus Get request failed: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("StarPlus can not parse body: " + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
	}
	body := string(b)
	//fmt.Println(body)
	if resp.StatusCode == 403 {
		if model.EnableLoger {
			Logger.Info("StarPlus request banned with status code: 403")
		}
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Request.URL.String())
	if resp.StatusCode == 200 {
		if resp.StatusCode == 302 || resp.Header.Get("Location") == "https://www.preview.starplus.com/unavailable" {
			if model.EnableLoger {
				Logger.Info("StarPlus region is inaccessible")
			}
			return model.Result{Name: name, Status: model.StatusNo}
		}
		region := utils.ReParse(body, `Region:\s+([A-Za-z]{2})`)
		if region != "" {
			loc := strings.ToLower(region)
			if utils.GetRegion(loc, model.StarPlusSupportCountry) {
				anotherCheck := AnotherStarPlus(c)
				result1, result2, result3 := utils.CheckDNS(hostname)
				unlockType := utils.GetUnlockType(result1, result2, result3)
				if anotherCheck.Err == nil && anotherCheck.Status == model.StatusYes {
					return model.Result{Name: name, Status: model.StatusYes, Region: loc, UnlockType: unlockType}
				} else {
					anotherCheck.Info = "Website: " + model.StatusYes
					anotherCheck.UnlockType = unlockType
					return anotherCheck
				}
			}
			return model.Result{Name: name, Status: model.StatusNo}
		}
	}
	if model.EnableLoger {
		Logger.Info(fmt.Sprintf("StarPlus unexpected response code: %d", resp.StatusCode))
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.starplus.com failed with code: %d", resp.StatusCode)}
}

// AnotherStarPlus
// StarPlus 的 另一个检测逻辑
func AnotherStarPlus(c *http.Client) model.Result {
	name := "Star+"
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	headers := map[string]string{
		"authorization": "c3RhciZicm93c2VyJjEuMC4w.COknIGCR7I6N0M5PGnlcdbESHGkNv7POwhFNL-_vIdg",
	}
	starcontent := "{\"query\":\"mutation registerDevice($input: RegisterDeviceInput!) " +
		"{\\n            registerDevice(registerDevice: $input) {\\n                grant " +
		"{\\n                    grantType\\n                    assertion\\n                " +
		"}\\n            }\\n        }\",\"variables\":{\"input\":{\"deviceFamily\":\"browser\"," +
		"\"applicationRuntime\":\"chrome\",\"deviceProfile\":\"windows\",\"deviceLanguage\":\"zh-CN\"," +
		"\"attributes\":{\"osDeviceIds\":[],\"manufacturer\":\"microsoft\",\"model\":null," +
		"\"operatingSystem\":\"windows\",\"operatingSystemVersion\":\"10.0\",\"browserName\":" +
		"\"chrome\",\"browserVersion\":\"96.0.4664\"}}}}"
	resp, body, err := utils.PostJson(c, "https://star.api.edge.bamgrid.com/graph/v1/device/graphql", starcontent, headers)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("star.api.edge.bamgrid.com post json err:" + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode >= 400 {
		if model.EnableLoger {
			Logger.Info("star.api.edge.bamgrid.com post json err with resp status code: " + resp.Status)
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("resp status code >= 400")}
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		if model.EnableLoger {
			Logger.Info("star.api.edge.bamgrid.com can not parse body:" + err.Error())
		}
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
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
	if model.EnableLoger {
		Logger.Info(fmt.Sprintf("www.starplus.com unexpected response code: %d", resp.StatusCode))
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.starplus.com another check failed with code: %d", resp.StatusCode)}
}
