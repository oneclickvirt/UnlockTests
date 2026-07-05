package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

// Claude 检测
func Claude(c *http.Client) model.Result {
	if c == nil {
		return model.Result{Name: "Claude"}
	}
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:             "Claude",
		hostname:         "claude.ai",
		url:              "https://claude.ai/",
		traceURL:         "https://claude.ai/cdn-cgi/trace",
		okCodes:          map[int]bool{http.StatusOK: true, http.StatusAccepted: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:          map[int]bool{http.StatusUnavailableForLegalReasons: true},
		bannedCodes:      map[int]bool{http.StatusForbidden: true},
		supportCountries: model.ClaudeSupportCountry,
		wafKeywords:      defaultAIWAFKeywords(),
	})
}
