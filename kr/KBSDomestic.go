package kr

import (
	"fmt"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// KBSDomestic
// vod.kbs.co.kr 仅 ipv4 且 get 请求
func KBSDomestic(request *gorequest.SuperAgent) model.Result {
	name := "KBS Domestic"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://vod.kbs.co.kr/index.html?source=episode&sname=vod&stype=vod&program_code=T2022-0690&program_id=PS-2022164275-01-000&broadcast_complete_yn=N&local_station_code=00&section_code=03"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	tempList := strings.Split(body, "\n")
	for _, line := range tempList {
		if strings.Contains(line, "ipck") && strings.Contains(line, "Domestic") {
			tpList := strings.Split(line, "Domestic")
			if strings.Contains(strings.Split(tpList[1], "\"")[1], "false") {
				return model.Result{Name: name, Status: model.StatusNo}
			} else if strings.Contains(strings.Split(tpList[1], "\"")[1], "true") {
				return model.Result{Name: name, Status: model.StatusYes}
			}
		}
	}
	if strings.Contains(body, "해당 영상은 저작권 등의 문제로") && strings.Contains(body, "서비스가 제공되지 않습니다") {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get vod.kbs.co.kr failed with code: %d", resp.StatusCode)}
}
