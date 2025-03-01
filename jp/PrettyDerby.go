package jp

import (
	"context"
	"net/http"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// PrettyDerby
// api-umamusume.cygames.jp 双栈 且 get 请求
func PrettyDerby(c *http.Client) model.Result {
	name := "Pretty Derby Japan"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://api-umamusume.cygames.jp/"
	headers := map[string]string{
		"User-Agent":                model.UA_Dalvik,
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Cache-Control":             "no-cache",
		"DNT":                       "1",
		"Pragma":                    "no-cache",
		"Sec-CH-UA":                 `"Not(A:Brand";v="99", "Microsoft Edge";v="133", "Chromium";v="133"` ,
		"Sec-CH-UA-Mobile":          "?0",
		"Sec-CH-UA-Platform":        "macOS",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
	}
	client := utils.ReqDefault(c)
	client = utils.SetReqHeaders(client, headers)
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.R().SetContext(ctx).Get(url)
		if err != nil {
			if err.Error() == `Get "https://api-umamusume.cygames.jp/": context deadline exceeded` {
				return model.Result{Name: name, Status: model.StatusNo}
			}
			continue // 发生错误时重试
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}