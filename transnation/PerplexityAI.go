package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func PerplexityAI(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Perplexity AI",
		hostname: "www.perplexity.ai",
		url:      "https://www.perplexity.ai/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}
