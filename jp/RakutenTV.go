package jp

import (
	"fmt"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// RakutenTV
// www.rakuten.tv 仅 ipv4 且 get 请求 带 cloudflare 的 5秒盾 无法使用 "is not available in your country"
// api.tv.rakuten.co.jp 仅 ipv4 且 get 请求 无盾可使用
func RakutenTV(request *gorequest.SuperAgent) model.Result {
	name := "Rakuten TV"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://api.tv.rakuten.co.jp/content/playinfo.json?content_id=476611&device_id=14&trailer=1&auth=0&log=0&serial_code=&tmp_eng_flag=1&multi_audio_support=1&_=1716694365356"
	resp, body, errs := request.Get(url).
		Set("connection", "keep-alive").
		Set("Cookie", "alt_id=kdPG3ErDszsWchi~f3P7Y3Mk; _ra=1716693934724|fbf06bf6-0e63-49bc-b5ae-ea8e785126ba; sec_token=6d518581124ba17c1b9968dca83aba7d441dcf88s%3A40%3A%220f817994db4925695da3375e3248a7552d981647%22%3B").
		Set("origin", "https://tv.rakuten.co.jp").
		Set("referer", "https://tv.rakuten.co.jp/").
		Timeout(10 * time.Second).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	// fmt.Println(body)
	if resp.StatusCode == 403 || strings.Contains(body, "海外からのアクセスのため、動画を再生できません。") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.rakuten.tv failed with code: %d", resp.StatusCode)}
}
