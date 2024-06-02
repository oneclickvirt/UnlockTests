package nl

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NPOStartPlus
// www.npo.nl 双栈 且 get 请求
func NPOStartPlus(request *gorequest.SuperAgent) model.Result {
	name := "NPO Start Plus"
	if request == nil {
		return model.Result{Name: name}
	}
	tokenURL := "https://www.npo.nl/start/api/domain/player-token?productId=LI_NL1_4188102"
	streamURL := "https://prod.npoplayer.nl/stream-link"
	referrerURL := "https://npo.nl/start/live?channel=NPO1"
	client := req.DefaultClient()
	client.ImpersonateChrome()
	resp, err := client.R().Get(tokenURL)
	defer resp.Body.Close()
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	// body := string(b)
	// fmt.Println(body)
	var res1 struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(b, &res1); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	token := res1.Token
	client.Headers.Set("Origin", "https://npo.nl")
	client.Headers.Set("Referer", "https://npo.nl/")
	client.Headers.Set("Content-Type", "application/json")
	client.Headers.Set("Authorization", token)
	resp2, err2 := client.R().SetBodyString(`{"profileName":"dash","drmType":"playready","referrerUrl":"` + referrerURL + `"}`).
		Post(streamURL)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp2.Body.Close()
	b, err = io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	// fmt.Println(body)
	// fmt.Println(resp2.StatusCode)
	// {"status":451,"body":"Dit programma mag niet bekeken worden vanaf jouw locatie."}
	if resp2.StatusCode == 451 || strings.Contains(body, "Dit programma mag niet bekeken worden vanaf jouw locatie.") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if resp2.StatusCode == 200 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else {
		return model.Result{Name: name, Status: model.StatusNo}
	}
}
