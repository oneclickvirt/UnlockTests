package transnation

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Instagram
// www.instagram.com 双栈 且 post 请求
func Instagram(c *http.Client) model.Result {
	name := "Instagram Licensed Audio"
	hostname := "www.instagram.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://www.instagram.com/graphql/query"
	payload := `av=0&__d=www&__user=0&__a=1&__req=6&__hs=20004.HYP%3Ainstagram_web_pkg.2.1..0.0&dpr=1&__ccg=UNKNOWN&__rev=1017147356&__s=atfa8i%3Afuzizb%3Agxoi6b&__hsi=7423427084978664422&__dyn=7xeUjG1mxu1syUbFp41twpUnwgU7SbzEdF8aUco2qwJw5ux609vCwjE1xoswaq0yE462mcw5Mx62G5UswoEcE7O2l0Fwqo31w9O1TwQzXwae4UaEW2G0AEco5G0zK5o4q0HUvw5rwSyES1TwVwDwHg2ZwrUdUbGwmk0zU8oC1Iwqo5q3e3zhA6bwIxe6V89F8uwm9EO2e2e0N9Wy8&__csr=l6g8QtdvveiiQsFl-cEGnKTVmPuvWjh8BJ4gyFr--VozoSHyQGjDBiKVUkWyfKmFlQl2AEZ2RBKLUyubAkwO49VUCquvqzoFzbK2ScogihqhrCK4qAgmAAqC-USibBDxaHJ5AyZoWezk8Kex2icy801dXo2CwhE7e0PU1pqz8Z0kGa1AwaS0O83pCxx166yet2H807JU1nFEaO03n8tx8gbwAg9U-0yU4VAKQ0Fy0acStpho9830oP8U8kGx5pEHo24CVo7i0DUdU66bwzweW15zgC9IMj5i0bTzE0aeE05ny05OE&__comet_req=7&lsd=AVpXrud7Bzo&jazoest=21036&__spin_r=1017147356&__spin_b=trunk&__spin_t=1728401306&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=PolarisPostActionLoadPostQueryQuery&variables=%7B%22shortcode%22%3A%22DAt7m0-P0u_%22%2C%22fetch_tagged_user_count%22%3Anull%2C%22hoisted_comment_id%22%3Anull%2C%22hoisted_reply_id%22%3Anull%7D&server_timestamps=true&doc_id=8845758582119845`
	headers := map[string]string{
		"Accept":                      "*/*",
		"Accept-Language":             "zh-CN,zh;q=0.9",
		"Connection":                  "keep-alive",
		"Content-Type":                "application/x-www-form-urlencoded",
		"Cookie":                      "csrftoken=IU6ZrkJP0ad06W4BvMKZfO; dpr=1.25; mid=ZwVGSwALAAGVQAfVcST--xXCdSEz; ig_did=ABFD3E31-FA87-4196-A81D-BC00BF6D22FA; ig_nrcb=1; ps_l=1; ps_n=1; datr=HEcFZ6tVauqNOeJjFQAXeJOE; wd=1439x1012",
		"Origin":                      "https://www.instagram.com",
		"Referer":                     "https://www.instagram.com/p/DAt7m0-P0u_/",
		"X-ASBD-ID":                   "129477",
		"X-CSRFToken":                 "IU6ZrkJP0ad06W4BvMKZfO",
		"X-FB-Friendly-Name":          "PolarisPostActionLoadPostQueryQuery",
		"X-FB-LSD":                    "AVpXrud7Bzo",
		"X-IG-App-ID":                 "936619743392459",
		"dpr":                         "1.25",
		"sec-ch-prefers-color-scheme": "light",
		"sec-ch-ua":                   `"Microsoft Edge";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`,
		"sec-ch-ua-full-version-list": `"Microsoft Edge";v="129.0.2792.65", "Not=A?Brand";v="8.0.0.0", "Chromium";v="129.0.6668.71"`,
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             `""`,
		"sec-ch-ua-platform":          `"Windows"`,
		"sec-ch-ua-platform-version":  `"10.0.0"`,
		"viewport-width":              "1640",
		"sec-fetch-dest":              "empty",
		"sec-fetch-mode":              "cors",
		"sec-fetch-site":              "same-origin",
		"User-Agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36 Edg/129.0.0.0",
	}

	resp, body, err := utils.PostJson(c, url, payload, headers)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	if resp.StatusCode == 200 {
		if strings.Contains(body, `"should_mute_audio":true`) {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	} else if resp.StatusCode == 429 {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Too Many Requests"}
	}
	return model.Result{
		Name:   name,
		Status: model.StatusUnexpected,
		Err:    fmt.Errorf("get www.instagram.com failed with code: %d", resp.StatusCode),
	}
}
