package eu

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
)

// RakutenTV
// gizmo.rakuten.tv 仅 ipv4 且 post 请求
func RakutenTV(request *gorequest.SuperAgent) model.Result {
	name := "Rakuten TV"
	url := "https://gizmo.rakuten.tv/v3/me/start?device_identifier=web&device_stream_audio_quality=2.0&device_stream_hdr_type=NONE&device_stream_video_quality=FHD"
	payload := `{"device_identifier":"web","device_metadata":{"app_version":"v5.5.22","audio_quality":"2.0","brand":"chrome","firmware":"XX.XX.XX","hdr":false,"model":"GENERIC","os":"Android OS","sdk":"112.0.0","serial_number":"not implemented","trusted_uid":false,"uid":"ab0dd3e8-5cae-4ad2-ba86-97af867e75c3","video_quality":"FHD","year":1970},"ifa_id":"b9c55e58-d5d0-41ed-becb-a54499731531"}`
	resp, body, errs := utils.PostJson(request, url, payload)
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	bodyString := string(body)
	if strings.Contains(bodyString, "forbidden_market") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "Not Available"}
	}
	if strings.Contains(bodyString, "forbidden_vpn") {
		return model.Result{Name: name, Status: model.StatusNo, Info: "VPN Forbidden"}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get gizmo.rakuten.tv failed with code: %d", resp.StatusCode)}
}
