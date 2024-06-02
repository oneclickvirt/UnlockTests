package tw

import (
	"encoding/json"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// LiTV
// www.litv.tv 仅 ipv4 且 post 请求
func LiTV(request *gorequest.SuperAgent) model.Result {
	name := "LiTV"
	if request == nil {
		return model.Result{Name: name}
	}
	headers := []map[string]string{
		{"Cookie": "PUID=34eb9a17-8834-4f83-855c-69382fd656fa; L_PUID=34eb9a17-8834-4f83-855c-69382fd656fa; device-id=f4d7faefc54f476bb2e7e27b7482469a"},
		{"Origin": "https://www.litv.tv"},
		{"Referer": "https://www.litv.tv/drama/watch/VOD00331042"},
		{"Priority": "u=1, i"},
	}
	resp, body, errs := utils.PostJson(request, "https://www.litv.tv/api/get-urls-no-auth",
		`{"AssetId": "vod71211-000001M001_1500K","MediaType": "vod","puid": "d66267c2-9c52-4b32-91b4-3e482943fe7e"}`,
		headers...,
	)
	if len(errs) > 0 {
		tp := AnotherLiTV()
		tp.Err = errs[0]
		return tp
	}
	bodyString := string(body)
	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "OutsideRegionError") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return AnotherLiTV()
}

// AnotherLiTV
// www.litv.tv 的另一个检测逻辑
func AnotherLiTV() model.Result {
	name := "LiTV"
	url := "https://www.litv.tv/vod/ajax/getUrl"
	payload := `{"type":"noauth","assetId":"vod44868-010001M001_800K","puid":"6bc49a81-aad2-425c-8124-5b16e9e01337"}`
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	request := gorequest.New()
	resp, body, errs := utils.PostJson(request, url, payload, headers)
	if errs != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	var jsonResponse map[string]interface{}
	err := json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	errorMessage, ok := jsonResponse["errorMessage"].(string)
	if !ok {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	switch errorMessage {
	case "null":
		return model.Result{Name: name, Status: model.StatusYes}
	case "vod.error.outsideregionerror":
		return model.Result{Name: name, Status: model.StatusNo}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
}
