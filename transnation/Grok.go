package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func Grok(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Grok",
		hostname: "grok.com",
		url:      "https://grok.com/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}
