package transnation

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

const dolaURL = "https://www.dola.com/chat/"

var dolaRegionRegex = regexp.MustCompile(`(?i)"inner_region"\s*:\s*"([A-Z]{2})"`)

func Dola(c *http.Client) model.Result {
	return checkDola(c, dolaURL, "www.dola.com")
}

func checkDola(c *http.Client, url, hostname string) model.Result {
	const name = "Dola AI"
	if c == nil {
		return model.Result{Name: name}
	}
	resp, err := utils.Req(c).R().
		SetHeader("User-Agent", model.UA_Browser).
		Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return model.Result{Name: name, Status: model.StatusRateLimited, Info: "HTTP 429"}
	case http.StatusForbidden:
		return model.Result{Name: name, Status: model.StatusBanned}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	match := dolaRegionRegex.FindSubmatch(body)
	if len(match) != 2 {
		return model.Result{Name: name, Status: model.StatusUnexpected}
	}
	return model.Result{Name: name, Status: model.StatusYes, Region: strings.ToLower(string(match[1]))}
}
