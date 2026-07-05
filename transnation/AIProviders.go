package transnation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func Poe(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:       "Poe",
		hostname:   "poe.com",
		url:        "https://poe.com/",
		noRedirect: true,
		okCodes:    map[int]bool{http.StatusOK: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:    map[int]bool{http.StatusForbidden: true},
	})
}

func PerplexityAI(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Perplexity AI",
		hostname: "www.perplexity.ai",
		url:      "https://www.perplexity.ai/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}

func MistralAI(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Mistral AI",
		hostname: "chat.mistral.ai",
		url:      "https://chat.mistral.ai/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}

func Grok(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Grok",
		hostname: "grok.com",
		url:      "https://grok.com/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}

func Coze(c *http.Client) model.Result {
	name := "Coze"
	hostname := "www.coze.com"
	if c == nil {
		return model.Result{Name: name}
	}
	resp, err := utils.Req(c).R().Get("https://www.coze.com/api/developer/get_login_info")
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return model.Result{Name: name, Status: model.StatusBanned}
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	var response struct {
		Code int `json:"code"`
		Data struct {
			IsForbiddenRegion bool   `json:"IsForbiddenRegion"`
			CountryCode       string `json:"CountryCode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &response); err != nil {
		body := string(b)
		if strings.Contains(body, "Your region is not supported") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if response.Data.IsForbiddenRegion {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if response.Data.CountryCode != "" {
		return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(response.Data.CountryCode)}
	}
	if response.Code == 0 {
		return model.Result{Name: name, Status: model.StatusYes}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected}
}

func DeepSeek(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "DeepSeek",
		hostname: "chat.deepseek.com",
		url:      "https://chat.deepseek.com/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}

func Kimi(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:       "Kimi",
		hostname:   "kimi.moonshot.cn",
		url:        "https://kimi.moonshot.cn/",
		noRedirect: true,
		okCodes: map[int]bool{
			http.StatusOK:               true,
			http.StatusMovedPermanently: true,
			http.StatusFound:            true,
		},
		noCodes: map[int]bool{http.StatusForbidden: true},
	})
}
