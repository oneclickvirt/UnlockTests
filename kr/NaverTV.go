package kr

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NaverTV
// apis.naver.com 仅 ipv4 且 get 请求
func NaverTV(request *gorequest.SuperAgent) model.Result {
	name := "Naver TV"
	if request == nil {
		return model.Result{Name: name}
	}
	ts := time.Now().UnixNano() / int64(time.Millisecond)
	baseURL := "https://apis.naver.com/"
	key := "nbxvs5nwNG9QKEWK0ADjYA4JZoujF4gHcIwvoCxFTPAeamq5eemvt5IWAYXxrbYM"
	signText := fmt.Sprintf("https://apis.naver.com/now_web2/now_web_api/v1/clips/31030608/play-info%d", ts)
	// 生成 HMAC-SHA1 签名并进行 base64 编码
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(signText))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// URL 对签名进行编码
	signatureEncoded := url.QueryEscape(signature)
	reqURL := fmt.Sprintf("%snow_web2/now_web_api/v1/clips/31030608/play-info?msgpad=%d&md=%s", baseURL, ts, signatureEncoded)
	// 进行请求
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(reqURL).
		Set("Host", "apis.naver.com").
		Set("Connection", "keep-alive").
		Set("Accept", "application/json, text/plain, */*").
		Set("User-Agent", "your-user-agent-here").
		Set("Origin", "https://tv.naver.com").
		Set("Referer", "https://tv.naver.com/v/31030608").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var res struct {
			Result struct {
				Play struct {
					Playable string `json:"playable"`
				} `json:"play"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(body), &res); err != nil {
			if strings.Contains(body, "NOT_COUNTRY_AVAILABLE") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			return model.Result{Name: name, Status: model.StatusErr, Err: err}
		}
		if res.Result.Play.Playable == "NOT_COUNTRY_AVAILABLE" {
			return model.Result{Name: name, Status: model.StatusNo}
		} else if res.Result.Play.Playable != "NOT_COUNTRY_AVAILABLE" {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get apis.naver.com failed with code: %d", resp.StatusCode)}
}
