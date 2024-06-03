package us

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// EncoreTVB
// edge.api.brightcove.com 仅 ipv4 且 get 请求
func EncoreTVB(request *gorequest.SuperAgent) model.Result {
	name := "EncoreTVB"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://edge.api.brightcove.com/playback/v1/accounts/5324042807001/videos/6005570109001"
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("User-Agent", model.UA_Browser)
	client.Headers.Set("Accept", "application/json;pk=BCpkADawqM2Gpjj8SlY2mj4FgJJMfUpxTNtHWXOItY1PvamzxGstJbsgc-zFOHkCVcKeeOhPUd9MNHEGJoVy1By1Hrlh9rOXArC5M5MTcChJGU6maC8qhQ4Y8W-QYtvi8Nq34bUb9IOvoKBLeNF4D9Avskfe9rtMoEjj6ImXu_i4oIhYS0dx7x1AgHvtAaZFFhq3LBGtR-ZcsSqxNzVg-4PRUI9zcytQkk_YJXndNSfhVdmYmnxkgx1XXisGv1FG5GOmEK4jZ_Ih0riX5icFnHrgniADr4bA2G7TYh4OeGBrYLyFN_BDOvq3nFGrXVWrTLhaYyjxOr4rZqJPKK2ybmMsq466Ke1ZtE-wNQ")
	client.Headers.Set("Origin", "https://www.encoretvb.com")
	resp, err := client.R().
		SetRetryCount(2).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second).Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	var res1 struct {
		ErrorSubcode string `json:"error_subcode"`
		AccountId    string `json:"account_id"`
	}
	var res2 []struct {
		ClientGeo    string `json:"client_geo"`
		ErrorSubcode string `json:"error_subcode"`
		ErrorCode    string `json:"error_code"`
		Message      string `json:"message"`
	}
	if err := json.Unmarshal(b, &res1); err != nil {
		if err := json.Unmarshal(b, &res2); err != nil {
			if strings.Contains(body, "CLIENT_GEO") || strings.Contains(body, "ACCESS_DENIED") {
				return model.Result{
					Name: name, Status: model.StatusNo,
				}
			}
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res2[0].ErrorSubcode == "CLIENT_GEO" {
			return model.Result{Name: name, Status: model.StatusNo, Region: res2[0].ClientGeo}
		}
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res1.AccountId != "0" {
		return model.Result{Name: name, Status: model.StatusYes, Region: "us"}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get edge.api.brightcove.com failed with code: %d", resp.StatusCode)}
}
