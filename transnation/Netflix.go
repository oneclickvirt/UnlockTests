package transnation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NetflixCDN
// api.fast.com 双栈 get 请求
func NetflixCDN(request *gorequest.SuperAgent) model.Result {
	name := "Netflix CDN"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api.fast.com/netflix/speedtest/v2?https=true&token=YXNkZmFzZGxmbnNkYWZoYXNkZmhrYWxm&urlCount=5"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
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
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res.Targets[0].Location.Country != "" {
		return model.Result{
			Name: name, Status: model.StatusYes,
			Region: res.Targets[0].Location.Country,
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get api.fast.com failed with code: %d", resp.StatusCode)}
}

// Netflix
// www.netflix.com 双栈 且 get 请求
func Netflix(request *gorequest.SuperAgent) model.Result {
	name := "Netflix"
	if request == nil {
		return model.Result{Name: name}
	}
	url1 := "https://www.netflix.com/title/81280792" // 乐高
	url2 := "https://www.netflix.com/title/70143836" // 绝命毒师
	url3 := "https://www.netflix.com/title/80018499" // Test Patterns
	request = request.Set("User-Agent", model.UA_Browser)
	resp1, _, errs1 := request.Get(url1).Retry(2, 5).End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp1.Body.Close()
	//if body1 == "" {
	//	return model.Result{
	//		Name: name, Status: model.StatusNo,
	//	}
	//}
	resp2, _, errs2 := request.Get(url2).Retry(2, 5).End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	defer resp2.Body.Close()
	//if body2 == "" {
	//	return model.Result{
	//		Name: name, Status: model.StatusNo,
	//	}
	//}
	if resp1.StatusCode == 404 && resp2.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusRestricted + " (Originals Only)"}
	}
	if resp1.StatusCode == 403 && resp2.StatusCode == 403 {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	if (resp1.StatusCode == 200 || resp1.StatusCode == 301) || (resp2.StatusCode == 200 || resp2.StatusCode == 301) {
		resp3, _, errs3 := request.Get(url3).Retry(2, 5).End()
		if len(errs3) > 0 {
			return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs3[0]}
		}
		defer resp3.Body.Close()
		//if body3 == "" {
		//	return model.Result{
		//		Name: name, Status: model.StatusNo,
		//	}
		//}
		u := resp3.Header.Get("location")
		if u == "" {
			return model.Result{Name: name, Status: model.StatusYes, Region: "us"}
		}
		//fmt.Println("nf", u)
		t := strings.SplitN(u, "/", 5)
		if len(t) < 5 {
			return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("can not find region")}
		}
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.SplitN(t[3], "-", 2)[0]}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.netflix.com failed with code: %d %d", resp1.StatusCode, resp2.StatusCode)}
}
