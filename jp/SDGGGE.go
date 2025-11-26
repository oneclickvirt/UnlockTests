package jp

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// SD Gundam G Generation Eternal
// 仅 ipv4 且 post 请求
func SDGGGE(c *http.Client) model.Result {
	name := "SD Gundam G"
	if c == nil {
		return model.Result{Name: name}
	}
	rawPayload := "1CR6PntuLeI3yaCYAZdOPxn18bOFYJxUiYtcavqqAHDCjc3C/wozplHYwfhykUStp7Bb/LAhV8aWQkS9sLliHCIgXBvDsWe4pwXvV3cSXkoaBfL23/zytEHlAatOi/32UVYLJhyUsegCRMMGREr2fXqyx970imQ35hqWVj/MRTHS9Bi8iqo9nIqSDTcQqVn3BbuyhJcz52nhfSda2may3QVHkH9QDdFjW9S/2re2cxE3iaE/DUbjB9H8KUpihQB1Emf88I0241ea7CAI1jHel6aZ5Ul4XjTf8ug3Rl/T80A="
	body, err := base64.StdEncoding.DecodeString(rawPayload)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	url := "https://api.gl.eternal.channel.or.jp/api/pvt/consent/view?user_id=649635267711712178"
	headers := map[string]string{
		"Host":                   "api.gl.eternal.channel.or.jp",
		"X-Content-Is-Encrypted": "True",
		"X-Language":             "hk",
		"Accept":                 "application/protobuf",
		"X-Unity-Version":        "2022.3.45f1",
		"X-Master-Url":           "https://clientdata.gl.eternal.channel.or.jp/prd-gl/catalogs/hr0phpfWDVahMJGQIk2OSd6hy35YpQZVKYAo6lKeld-9scMGJw2KTnBDGbS04Gw-i25avFTH55K-yU9TCX2OkQ.json",
		"X-Language-Master-Url":  "https://clientdata.gl.eternal.channel.or.jp/prd-gl/language_catalogs/hk/F-HORjFKHLai8nLXUdPyQRqzexZNPKIn2O36Hgd2Bxm2RysBNS0-PQHQwfHXEOONog0w5yULtewBaVk-Ndf6nQ.json",
		"x-app-version-hash":     "20928",
		"x-token":                "e5df59f1-8588-4477-a887-5fe854895493Mj0jmtfbgIhQOUmHQE1W7sLq7G5eSBqcFWqldSPjy6s=",
		"Accept-Language":        "zh-CN,zh-Hans;q=0.9",
		"User-Agent":             "GETERNAL/25041500 CFNetwork/3826.400.120 Darwin/24.3.0",
		"Connection":             "keep-alive",
		"Content-Type":           "application/protobuf",
	}
	client := utils.ReqDefault(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().SetBody(body).Post(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return model.Result{Name: name, Status: model.StatusYes}
	case 483:
		return model.Result{Name: name, Status: model.StatusNo}
	default:
		return model.Result{
			Name:   name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("SDGGGE unknown with code: %d", resp.StatusCode),
		}
	}
}
