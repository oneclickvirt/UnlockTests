package us

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// NBCTV
// geolocation.digitalsvc.apps.nbcuni.com 双栈 get 请求
func NBCTV(request *gorequest.SuperAgent) model.Result {
	name := "NBC TV"
	if request == nil {
		return model.Result{Name: name}
	}
	fakeUuid, _ := uuid.NewV4()
	url := "https://geolocation.digitalsvc.apps.nbcuni.com/geolocation/live/usa"
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Headers.Set("accept-language", "en-US,en;q=0.9")
	client.Headers.Set("app-session-id", fakeUuid.String())
	client.Headers.Set("authorization", "NBC-Basic key=\"usa_live\", version=\"3.0\", type=\"cpc\"")
	client.Headers.Set("client", "oneapp")
	client.Headers.Set("content-type", "application/json")
	client.Headers.Set("origin", "https://www.nbc.com")
	client.Headers.Set("referer", "https://www.nbc.com/")
	client.Headers.Set("sec-ch-ua", model.UA_SecCHUA)
	client.Headers.Set("sec-ch-ua-mobile", "?0")
	client.Headers.Set("sec-ch-ua-platform", "\"Windows\"")
	client.Headers.Set("sec-fetch-dest", "empty")
	client.Headers.Set("sec-fetch-mode", "cors")
	client.Headers.Set("sec-fetch-site", "cross-site")
	resp, err := client.R().
		SetRetryCount(2).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second).
		SetBodyJsonString(`{"adobeMvpdId":null,"serviceZip":null,"device":"web"}`).
		Post(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	// fmt.Println(body)
	if strings.Contains(body, `"restricted":false`) {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if strings.Contains(body, `"restricted":true`) || body == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("getgeolocation.digitalsvc.apps.nbcuni.com failed with code: %d", resp.StatusCode)}
}
