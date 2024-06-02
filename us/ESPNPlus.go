package us

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ESPNPlus
// espn.api.edge.bamgrid.com 双栈 且 post 请求 可能 有 cloudflare 的5秒盾
func ESPNPlus(request *gorequest.SuperAgent) model.Result {
	name := "ESPN+"
	url1 := "https://espn.api.edge.bamgrid.com/token"
	data1 := `grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Atoken-exchange&latitude=0&longitude=0&platform=browser&subject_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJjYWJmMDNkMi0xMmEyLTQ0YjYtODJjOS1lOWJkZGNhMzYwNjkiLCJhdWQiOiJ1cm46YmFtdGVjaDpzZXJ2aWNlOnRva2VuIiwibmJmIjoxNjMyMjMwMTY4LCJpc3MiOiJ1cm46YmFtdGVjaDpzZXJ2aWNlOmRldmljZSIsImV4cCI6MjQ5NjIzMDE2OCwiaWF0IjoxNjMyMjMwMTY4LCJqdGkiOiJhYTI0ZWI5Yi1kNWM4LTQ5ODctYWI4ZS1jMDdhMWVhMDgxNzAifQ.8RQ-44KqmctKgdXdQ7E1DmmWYq0gIZsQw3vRL8RvCtrM_hSEHa-CkTGIFpSLpJw8sMlmTUp5ZGwvhghX-4HXfg&subject_token_type=urn%3Abamtech%3Aparams%3Aoauth%3Atoken-type%3Adevice`
	resp, body, errs := request.Post(url1).Type("form").
		Send(data1).
		Set("authorization", "Bearer ZXNwbiZicm93c2VyJjEuMC4w.ptUt7QxsteaRruuPmGZFaJByOoqKvDP2a5YkInHrc7c").End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// fmt.Println(body)
	if strings.Contains(body, "forbidden-location") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	var res struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}

	url2 := "https://espn.api.edge.bamgrid.com/graph/v1/device/graphql"
	data2 := `{"query":"mutation registerDevice($input: RegisterDeviceInput!) {\n            registerDevice(registerDevice: $input) {\n                grant {\n                    grantType\n                    assertion\n                }\n            }\n        }","variables":{"input":{"deviceFamily":"browser","applicationRuntime":"chrome","deviceProfile":"windows","deviceLanguage":"zh-CN","attributes":{"osDeviceIds":[],"manufacturer":"microsoft","model":null,"operatingSystem":"windows","operatingSystemVersion":"10.0","browserName":"chrome","browserVersion":"96.0.4664"}}}}`
	resp2, body2, errs2 := request.Post(url2).Type("json").
		Send(data2).
		Set("authorization", "ZXNwbiZicm93c2VyJjEuMC4w.ptUt7QxsteaRruuPmGZFaJByOoqKvDP2a5YkInHrc7c").End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	var res2 struct {
		Extensions struct {
			Sdk struct {
				Session struct {
					Location struct {
						CountryCode string `json:"countryCode"`
					}
					InSupportedLocation bool `json:"inSupportedLocation"`
				}
			}
		}
	}
	if err := json.Unmarshal([]byte(body2), &res2); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if res2.Extensions.Sdk.Session.Location.CountryCode == "US" && res2.Extensions.Sdk.Session.InSupportedLocation {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get espn.api.edge.bamgrid.com failed with code: %d", resp.StatusCode)}
}
