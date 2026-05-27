package transnation

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Spotify
// open.spotify.com 检测市场可用性，避免 spclient 的代理检测
func Spotify(c *http.Client) model.Result {
	name := "Spotify Registration"
	hostname := "open.spotify.com"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://open.spotify.com/"
	headers := map[string]string{
		"User-Agent":      model.UA_Browser,
		"Accept-Language": "en",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	body := string(b)
	// Extract the appServerConfig base64 JSON embedded in the page
	re := regexp.MustCompile(`<script[^>]+id="appServerConfig"[^>]*type="text/plain"[^>]*>([^<]+)</script>`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(matches[1]))
	if err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("base64 decode failed: %w", err)}
	}
	var config struct {
		Market string `json:"market"`
	}
	if err := json.Unmarshal(decoded, &config); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if config.Market == "" {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	result1, result2, result3 := utils.CheckDNS(hostname)
	unlockType := utils.GetUnlockType(result1, result2, result3)
	return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType,
		Region: strings.ToLower(config.Market)}
}
