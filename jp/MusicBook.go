package jp

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// MusicBook
// overseaauth.music-book.jp 仅 ipv4 且 get 请求
func MusicBook(request *gorequest.SuperAgent) model.Result {
	name := "music.jp"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://overseaauth.music-book.jp/globalIpcheck.js"
	request = request.Set("User-Agent", model.UA_Browser)
	resp, body, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	if resp.StatusCode == 403 || resp.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 && body != "" {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get overseaauth.music-book.jp failed with code: %d", resp.StatusCode)}
}
