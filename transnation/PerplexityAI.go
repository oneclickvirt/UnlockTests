package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func PerplexityAI(c *http.Client) model.Result {
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:                "Perplexity AI",
		hostname:            "www.perplexity.ai",
		url:                 "https://www.perplexity.ai/",
		traceURL:            "https://www.perplexity.ai/cdn-cgi/trace",
		okCodes:             map[int]bool{http.StatusOK: true, http.StatusAccepted: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:             map[int]bool{http.StatusUnavailableForLegalReasons: true},
		forbiddenCodes:      map[int]bool{http.StatusForbidden: true},
		restrictedCountries: aiGlobalRestrictedCountries,
		wafKeywords:         defaultAIWAFKeywords(),
	})
}
