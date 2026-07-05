package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func Poe(c *http.Client) model.Result {
	return checkAIRegionalStatus(c, aiRegionalProbe{
		name:        "Poe",
		hostname:    "poe.com",
		url:         "https://poe.com/",
		traceURL:    "https://poe.com/cdn-cgi/trace",
		noRedirect:  true,
		okCodes:     map[int]bool{http.StatusOK: true, http.StatusMovedPermanently: true, http.StatusFound: true, http.StatusTemporaryRedirect: true, http.StatusPermanentRedirect: true},
		noCodes:     map[int]bool{http.StatusForbidden: true, http.StatusUnavailableForLegalReasons: true},
		wafKeywords: defaultAIWAFKeywords(),
	})
}
