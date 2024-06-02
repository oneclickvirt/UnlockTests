package jp

import (
	"fmt"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// Mora
// mora.jp 仅 ipv4 且 get 请求
func Mora(request *gorequest.SuperAgent) model.Result {
	name := "Mora"
	if request == nil {
		return model.Result{Name: name}
	}
	url := "https://mora.jp/buy?__requestToken=1713764407153&returnUrl=https%3A%2F%2Fmora.jp%2Fpackage%2F43000087%2FTFDS01006B00Z%2F%3Ffmid%3DTOPRNKS%26trackMaterialNo%3D31168909&fromMoraUx=false&deleteMaterial="
	request = request.Set("User-Agent", model.UA_Browser).
		Set("host", "mora.jp")
	resp, _, errs := request.Get(url).Retry(2, 5).End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 403 || resp.StatusCode == 451 || resp.StatusCode == 500 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get mora.jp failed with code: %d", resp.StatusCode)}
}
