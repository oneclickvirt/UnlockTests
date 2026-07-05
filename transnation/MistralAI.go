package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func MistralAI(c *http.Client) model.Result {
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:                "Mistral AI",
		hostname:            "chat.mistral.ai",
		url:                 "https://chat.mistral.ai/",
		traceURL:            "https://chat.mistral.ai/cdn-cgi/trace",
		okCodes:             map[int]bool{http.StatusOK: true, http.StatusAccepted: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:             map[int]bool{http.StatusForbidden: true, http.StatusUnavailableForLegalReasons: true},
		restrictedCountries: mistralAIRestrictedCountries,
		wafKeywords:         defaultAIWAFKeywords(),
	})
}
