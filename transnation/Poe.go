package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

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
