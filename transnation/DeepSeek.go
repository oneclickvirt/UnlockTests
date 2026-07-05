package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func DeepSeek(c *http.Client) model.Result {
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:        "DeepSeek",
		hostname:    "chat.deepseek.com",
		url:         "https://chat.deepseek.com/",
		traceURL:    "https://chat.deepseek.com/cdn-cgi/trace",
		okCodes:     map[int]bool{http.StatusOK: true, http.StatusAccepted: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:     map[int]bool{http.StatusUnavailableForLegalReasons: true},
		bannedCodes: map[int]bool{http.StatusForbidden: true},
		wafKeywords: defaultAIWAFKeywords(),
	})
}
