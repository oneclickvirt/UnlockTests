package transnation

import (
	"fmt"
	"net/http"

	req "github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

type aiStatusProbe struct {
	name       string
	hostname   string
	url        string
	noRedirect bool
	okCodes    map[int]bool
	noCodes    map[int]bool
}

func checkAIStatus(c *http.Client, probe aiStatusProbe) model.Result {
	if c == nil {
		return model.Result{Name: probe.name}
	}
	client := utils.Req(c)
	if probe.noRedirect {
		client.SetRedirectPolicy(req.NoRedirectPolicy())
	}
	resp, err := client.R().Get(probe.url)
	if err != nil {
		return utils.HandleNetworkError(c, probe.hostname, err, probe.name)
	}
	defer resp.Body.Close()
	switch {
	case probe.okCodes[resp.StatusCode]:
		return model.Result{Name: probe.name, Status: model.StatusYes}
	case probe.noCodes[resp.StatusCode]:
		return model.Result{Name: probe.name, Status: model.StatusNo}
	default:
		return model.Result{
			Name:   probe.name,
			Status: model.StatusUnexpected,
			Err:    fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
	}
}
