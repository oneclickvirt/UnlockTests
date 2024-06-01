package nl

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NPOStartPlus
// www.npo.nl 双栈 且 get 请求
func NPOStartPlus(request *gorequest.SuperAgent) model.Result {
	name := "NPO Start Plus"
	tokenURL := "https://www.npo.nl/start/api/domain/player-token?productId=LI_NL1_4188102"
	streamURL := "https://prod.npoplayer.nl/stream-link"
	referrerURL := "https://npo.nl/start/live?channel=NPO1"
	resp, body, errs1 := request.Get(tokenURL).
		Set("User-Agent", model.UA_Browser).
		Set("Host", "www.npo.nl").
		Set("Connection", "keep-alive").
		Set("Accept", "application/json, text/plain, */*").
		Set("Referer", referrerURL).
		End()
	if len(errs1) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs1[0]}
	}
	defer resp.Body.Close()
	token := body
	resp2, _, errs2 := request.Post(streamURL).
		Set("User-Agent", model.UA_Browser).
		Set("Accept", "*/*").
		Set("Authorization", token).
		Set("Content-Type", "application/json").
		Set("Origin", "https://npo.nl").
		Set("Referer", "https://npo.nl/").
		Send(`{"profileName":"dash","drmType":"playready","referrerUrl":"` + referrerURL + `"}`).
		End()
	if len(errs2) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
	}
	if resp2.StatusCode == 451 {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp2.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
