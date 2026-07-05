package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func Grok(c *http.Client) model.Result {
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:                "Grok",
		hostname:            "grok.com",
		url:                 "https://grok.com/",
		traceURL:            "https://grok.com/cdn-cgi/trace",
		okCodes:             map[int]bool{http.StatusOK: true, http.StatusAccepted: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		forbiddenCodes:      map[int]bool{http.StatusForbidden: true},
		restrictedCountries: aiGlobalRestrictedCountries,
		wafKeywords:         defaultAIWAFKeywords(),
	})
}
